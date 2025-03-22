package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
)

func (a *Analyzer) checkScope(scp *nodes.ScopeNode) []string {
	errs := make([]string, 0)
	// check if datatypes exist and var decls match type
	for _, n := range scp.Children {
		switch v := n.(type) {
		case *nodes.VariableDeclarationNode:
			{
				for _, tn := range v.Declarations {
					if op := tn.(*nodes.InfixOperatorNode); op.Op == lexer.ASSN {
						varname := op.Left.(*nodes.IdentifierNode)
						a.definedSyms[varname.Name()] = v.DataT
					} else {
						errs = append(errs, fmt.Sprintf("Syntax error in variable declaration at line <TODO>: %s not allowed", op.Op))

					}
				}
				// a.definedSyms[v.]
			}
		}
	}
	return []string{}
}
