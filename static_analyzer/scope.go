package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
)

func (a *Analyzer) checkScope(scp *nodes.ScopeNode, returnType nodes.DataType) {
	// check if datatypes exist and var decls match type
	for _, n := range scp.Children {
		switch v := n.(type) {
		case *nodes.VariableDeclarationNode:
			{
				for _, tn := range v.Declarations {
					if op := tn.(*nodes.InfixOperatorNode); op.Op == lexer.ASSN {
						varname := op.Left.(*nodes.IdentifierNode)
						// todo: Check if v.DataT itself is valid
						a.definedSyms[varname.Name()] = v.DataT
						// rval should have same type
						rvalType := computeType(op.Right, a)
						if !rvalType.Equals(v.DataT) {
							a.AddError(
								tn.Range().Start,
								utils.TypeError,
								fmt.Sprintf("Cannot assign %s to variable of type %s", utils.Cyan(rvalType.Text()), utils.Cyan(v.DataT.Text())),
							)
						}
					} else {
						a.AddError(
							tn.Range().Start,
							utils.SyntaxError,
							fmt.Sprintf("%s not allowed", utils.Red(op.Op)),
						)
					}
				}
			}
		case *nodes.ReturnNode:
			{
				if ct := computeType(v.Value, a); !ct.Equals(returnType) {
					a.AddError(v.Range().Start, utils.TypeError,
						fmt.Sprintf("Expected to return value of type %s but found %s", utils.Cyan(returnType.Text()), utils.Cyan(ct.Text())))
				}
			}
		case *nodes.FuncCallNode:
			{
				funcType := computeType(v.Callee, a)
				ftyp, ok := funcType.(*nodes.FuncType)
				if !ok {
					a.AddError(
						v.Range().Start,
						utils.TypeError,
						fmt.Sprintf("Type is not callable: %s", utils.Cyan(funcType.Text())),
					)
					return
				}

				if len(v.Args) != len(ftyp.ArgTypes) {
					a.AddError(
						v.Range().Start,
						utils.TypeError,
						fmt.Sprintf("Function %s expects %s parameters, but supplied %s", utils.Blue(v.Callee.String("")), utils.Yellow(fmt.Sprint(len(ftyp.ArgTypes))), utils.Yellow(fmt.Sprint(len(v.Args)))),
					)
					return
				}
				for i, k := range v.Args {
					expT := ftyp.ArgTypes[i]
					passedT := computeType(k, a)

					if !expT.Equals(passedT) {
						a.AddError(
							v.Range().Start,
							utils.TypeError,
							fmt.Sprintf("%d th parameter to function %s should be of type %s, not %s", i, utils.Blue(v.Callee.String("")), utils.Cyan(expT.Text()), utils.Cyan(passedT.Text())),
						)
					}
				}

			}
		default:
			a.AddError(
				v.Range().Start,
				utils.UndefinedError,
				fmt.Sprintf("Can't check for %T", v),
			)
		}
	}
	return
}
