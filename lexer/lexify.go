package lexer

import (
	"bytes"
	"os"
	"strings"
	g "toylingo/globals"
	"toylingo/utils"
)

func Lexify(path string) *Node {

	filecontent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	//replace comments with spaces
	for i := 0; i < len(filecontent); i++ {
		if string(filecontent[i:i+2]) == "//" {
			for j := i; j < len(filecontent); j++ {
				if filecontent[j] == '\n' {
					i = j
					break
				}
				filecontent[j] = ' '
			}
		}
	}

	//decompose template literals

	// fmt.Println(string(filecontent))
	stringliterals := make([]string, 0)
	//placeholder for strings
	for i := 0; i < len(filecontent); i++ {
		if filecontent[i] == '`' {
			for j := i + 1; j < len(filecontent); j++ {
				if filecontent[j] == '`' {
					str := string(filecontent[i : j+1])

					stringliterals = append(stringliterals, str)
					i = j + 1
					break
				}
			}
		}
	}
	//parse escape sequence from all string literals
	// for i := 0; i < len(stringliterals); i++ {
	// 	stringliterals[i] = utils.ParseEscapeSequence(stringliterals[i])

	// }
	for i := 0; i < len(stringliterals); i++ {
		filecontent = bytes.ReplaceAll(filecontent, []byte(stringliterals[i]), []byte(" __STR__ "))
	}

	toPad := [...]string{"{", "}", ";", ":", "(", ")", ".", "=", "*", "/", "+", "-", "<", ">", "!", "|", "&", ","}
	for i := 0; i < len(toPad); i++ {
		filecontent = bytes.ReplaceAll(filecontent, []byte(toPad[i]), []byte(" "+toPad[i]+" "))
	}
	filecontent = append(filecontent, []byte(" ")...)

	// fmt.Println(string(filecontent))
	// fmt.Println(stringliterals)
	tokens := &Node{TokenType{"start", ""}, nil}
	ret := tokens
	temp := ""
	for i := 0; i < len(filecontent); i++ {
		c := filecontent[i]
		if c == ' ' || c == '\n' || c == '\t' {
			// fmt.Println("word ", temp)
			if addToken(strings.Trim(temp, " "), tokens) {
				tokens = tokens.Next
			}
			temp = ""
			continue
		}
		temp += string(c)

	}

	//coalesce decimal nos (4.67 etc) into one
	for node := ret; node != nil; node = node.Next {
		if node.Val.Type == "NUMBER" && node.Next != nil && node.Next.Val.Type == "DOT" && node.Next.Next != nil && node.Next.Next.Val.Type == "NUMBER" {
			node.Val.Ref = node.Val.Ref + "." + node.Next.Next.Val.Ref
			node.Next = node.Next.Next.Next
		}
	}

	//coalesce multicharacter operators into one
	for node := ret; node != nil; node = node.Next {
		if (utils.IsOneOf(node.Val.Ref, "<>=!") && node.Next != nil && node.Next.Val.Ref == "=") || (node.Val.Ref == "|" && node.Next != nil && node.Next.Val.Ref == "|") || (node.Val.Ref == "&" && node.Next != nil && node.Next.Val.Ref == "&") {
			node.Val.Ref = node.Val.Ref + node.Next.Val.Ref
			node.Next = node.Next.Next
		}
	}

	//replace placeholder for strings
	count := 0
	for node := ret; node != nil; node = node.Next {
		if node.Val.Type == "IDENTIFIER" && node.Val.Ref == "__STR__" {
			node.Val.Ref = utils.ParseEscapeSequence(stringliterals[count])
			count++
			node.Val.Type = "STRING_LITERAL"
		}
	}

	return ret
}

func addToken(temp string, tokens *Node) bool {
	if temp == " " || temp == "\n" || temp == "\t" || temp == "" || temp == "\r" {
		return false
	}
	switch strings.Trim(temp, " ") {
	case g.SCOPE_START:
		tokens.Next = &Node{TokenType{"SCOPE_START", g.SCOPE_START}, nil}
	case g.SCOPE_END:
		tokens.Next = &Node{TokenType{"SCOPE_END", g.SCOPE_END}, nil}
	case g.OPEN_PAREN:
		tokens.Next = &Node{TokenType{"OPEN_PAREN", g.OPEN_PAREN}, nil}
	case g.CLOSE_PAREN:
		tokens.Next = &Node{TokenType{"CLOSE_PAREN", g.CLOSE_PAREN}, nil}
	case g.COLON:
		tokens.Next = &Node{TokenType{"COLON", g.COLON}, nil}
	case g.SEMICOLON:
		tokens.Next = &Node{TokenType{"SEMICOLON", g.SEMICOLON}, nil}
	case g.LET:
		tokens.Next = &Node{TokenType{"LET", g.LET}, nil}
	case g.IF:
		tokens.Next = &Node{TokenType{"IF", g.IF}, nil}
	case g.ELSE_IF:
		tokens.Next = &Node{TokenType{"ELSE IF", g.ELSE_IF}, nil}
	case g.ELSE:
		tokens.Next = &Node{TokenType{"ELSE", g.ELSE}, nil}
	case g.LOOP:
		tokens.Next = &Node{TokenType{"LOOP", g.LOOP}, nil}
	case g.BREAK:
		tokens.Next = &Node{TokenType{"BREAK", g.BREAK}, nil}
	case g.DOT:
		tokens.Next = &Node{TokenType{"DOT", g.DOT}, nil}
	case g.EQUALS:
		tokens.Next = &Node{TokenType{"COMMA", g.EQUALS}, nil}
	case g.COMMA:
		tokens.Next = &Node{TokenType{"FUNCTION", g.COMMA}, nil}
	case g.RETURN:
		tokens.Next = &Node{TokenType{"RETURN", g.RETURN}, nil}
	case g.FUNCTION:
		tokens.Next = &Node{TokenType{"FUNCTION", g.FUNCTION}, nil}

	default:
		if utils.IsNumber(temp) {
			tokens.Next = &Node{TokenType{"NUMBER", temp}, nil}
		} else if isOperator(temp) {
			tokens.Next = &Node{TokenType{"OPERATOR", temp}, nil}
		} else {
			tokens.Next = &Node{TokenType{"IDENTIFIER", temp}, nil}
		}

	}
	return true
}

func isOperator(temp string) bool {
	operators := "=+-*/<>#!|&"
	return utils.IsOneOf(temp, operators)
}
