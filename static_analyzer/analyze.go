package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	// "he++/utils"
)

/**
Things to check for:
- idents not used before decl, no duplicate decl
- no expressions (OPERATOR) in top level scope
- no break, continue in non loop scope
- all code paths should have return
- type inference in var decl
- type consistency in func return
- normalized types
*/

type Analyzer struct {
	functionDecls map[string]*nodes.FuncNode
	definedTypes  map[string]DataTypeInfo
	// symname: normalized_typename
	definedSyms map[string]nodes.DataType
}

func MakeAnalyzer() Analyzer {

	return Analyzer{
		functionDecls: make(map[string]*nodes.FuncNode),
		definedTypes:  getPrimitiveTypeDefns(),
		definedSyms: make(map[string]nodes.DataType),
	}
}

func (a *Analyzer) AnalyzeAST(n *nodes.SourceFileNode) []string {
	errs := make([]string, 0)
	for _, ch := range n.Children {
		if ch.Type() == nodes.FUNCTION {
			funcNode := ch.(*nodes.FuncNode)
			a.registerFunctionDecl(funcNode)
		}
	}

	for _, ch := range n.Children {
		if ch.Type() == nodes.FUNCTION {
			funcNode := ch.(*nodes.FuncNode)
			errs = append(errs, a.checkFunctionDef(funcNode)...)
			fmt.Println(funcNode.Name, funcNode.ReturnType)
		}
	}

	return errs
}
