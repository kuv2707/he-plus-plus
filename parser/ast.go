package parser

import (
	"fmt"
	"he++/globals"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"sort"
	"strings"
)

func (p *Parser) ParseAST() *nodes.SourceFileNode {
	if !p.tokenStream.HasTokens() {
		parsingError("No tokens to parse!", -1)
		return nil
	}
	root := nodes.MakeSourceFileNode()
	parseStatements(p, root)

	// todo: add some postprocessing
	sort.SliceStable(root.Children, func(i, j int) bool {
		return root.Children[i].Type() == nodes.VAR_DECL
	})
	return root
}

func parseScope(p *Parser) *nodes.ScopeNode {
	t := p.tokenStream
	t.ConsumeOnlyIf(lexer.LPAREN)
	scopeNode := nodes.MakeScopeNode()
	parseStatements(p, scopeNode)
	t.ConsumeOnlyIf(lexer.RPAREN)
	return scopeNode
}

func parseStatements(p *Parser, scope nodes.StatementsContainer) nodes.StatementsContainer {
OUT:
	for p.tokenStream.HasTokens() {
		curr := p.tokenStream.Current()
		if curr.Type() == lexer.BRACKET && curr.Text() == lexer.RPAREN {
			break OUT
		}
		parselet, exists := p.scopeParselets[curr.Text()]
		fmt.Println(curr, exists)
		if !exists {
			expr := parseExpression(p, 0)
			if expr != nil {
				scope.AddChild(expr)
			} else {
				parsingError(fmt.Sprintf("Cannot parse %s", globals.Red(curr.Text())), curr.LineNo())
			}
			if !p.tokenStream.HasTokens() {
				break OUT
			}
		} else {
			scope.AddChild(parselet(p))
		}
	}
	return scope
}

func parseFunction(p *Parser) nodes.TreeNode {
	t := p.tokenStream
	t.ConsumeOnlyIf(lexer.FUNCTION)
	funcName := t.Consume()
	funcNode := nodes.MakeFunctionNode(funcName.Text())
	t.ConsumeOnlyIf(lexer.OPEN_PAREN)
	for t.Current().Text() != lexer.CLOSE_PAREN {
		varName := t.Consume()
		dataType := t.Consume()
		funcNode.AddArg(varName.Text(), nodes.DataType{Text:dataType.Text()})
		t.ConsumeIf(lexer.COMMA)
	}
	t.ConsumeOnlyIf(lexer.CLOSE_PAREN)
	funcNode.ReturnType = t.Consume().Text()
	funcNode.Scope = parseScope(p)
	return funcNode
}

func parseReturnStatement(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.RETURN)
	return nodes.MakeReturnNode(parseExpression(p, 0))
}

func parseDataType(p *Parser) nodes.DataType {
	compos := make([]string,0)
	for p.tokenStream.Current().Text() != "=" {
		compos = append(compos, p.tokenStream.Consume().Text())
	}
	p.tokenStream.Unread(1)
	return nodes.DataType{
		Text:strings.Join(compos[0 : len(compos)-1], ""),
	} 
}

func parseVariableDeclaration(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.LET)
	varDec := nodes.MakeVariableDeclarationNode()
	// todo: add type
	varDec.SetDataType(parseDataType(p))
	varDec.AddDeclaration(parseExpression(p, 0))
	for p.tokenStream.Current().Text() == lexer.COMMA {
		p.tokenStream.Consume()
		varDec.AddDeclaration(parseExpression(p, 0))
	}
	return varDec
}

func parseIfStatement(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.IF)
	ifCond := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.THEN)
	ifScope := parseScope(p)
	p.tokenStream.ConsumeOnlyIf(lexer.ELSE)
	if p.tokenStream.Current().Text() == lexer.IF {
		elseScope := parseIfStatement(p)
		return nodes.MakeIfNode(ifCond, ifScope, elseScope.(*nodes.IfNode))
	}
	elseScope := parseScope(p)
	return nodes.MakeIfNode(ifCond, ifScope, elseScope)
}

func parseLoopStatement(p *Parser) nodes.TreeNode {
	if p.tokenStream.Current().Text() == lexer.FOR {
		return parseForLoop(p)
	} else {
		return parseWhileLoop(p)
	}

}

func parseWhileLoop(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.WHILE)
	p.tokenStream.ConsumeOnlyIf(lexer.THAT)
	condNode := parseExpression(p, 0)
	return nodes.MakeLoopNode(&nodes.EmptyPlaceholderNode{}, condNode, &nodes.EmptyPlaceholderNode{}, parseScope(p))
}

func parseForLoop(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.FOR)
	varDecl := parseVariableDeclaration(p)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	condNode := parseExpression(p, 0)
	p.tokenStream.ConsumeOnlyIf(lexer.SEMICOLON)
	updNode := parseExpression(p, 0)
	return nodes.MakeLoopNode(varDecl, condNode, updNode, parseScope(p))
}

func parseStruct(p *Parser) nodes.TreeNode {
	p.tokenStream.ConsumeOnlyIf(lexer.STRUCT)
	var name string
	if p.tokenStream.Current().Text() == lexer.LPAREN {
		// anonymous struct
		name = "" // todo: give some label
	} else {
		name = p.tokenStream.Current().Text()
		p.tokenStream.ConsumeIfType(lexer.IDENTIFIER)
	}
	p.tokenStream.ConsumeOnlyIf(lexer.LPAREN)
	strct := nodes.StructNode{Name: name, Fields: make(map[string]nodes.StructField)}
	for p.tokenStream.Current().Text() != lexer.RPAREN {
		sname := p.tokenStream.Consume().Text()
		fmt.Print("-->", sname)
		stype := getStructType(p)
		strct.Fields[sname] = nodes.StructField{Name: sname, FieldType: stype}
	}
	p.tokenStream.ConsumeOnlyIf(lexer.RPAREN)
	return strct
}

func getStructType(p *Parser) nodes.TreeNode {
	if p.tokenStream.Current().Text() == lexer.STRUCT {
		// embedded struct
		return parseStruct(p)
	} else {
		return parseIdentifier(p)
	}
}
