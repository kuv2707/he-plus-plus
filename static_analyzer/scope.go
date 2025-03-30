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
						// rval should have same type
						rvalType := computeType(op.Right, a)
						if !rvalType.Equals(v.DataT) {
							errs = append(errs, fmt.Sprintf("Cannot assign %s to variable of type %s", rvalType.Text(), v.DataT.Text()))
						}
					} else {
						errs = append(errs, fmt.Sprintf("Syntax error in variable declaration at line <TODO>: %s not allowed", op.Op))
					}
				}
			}
		case *nodes.ReturnNode:
			{
				// no op
			}
		case *nodes.FuncCallNode:
			{
				// fmt.Printf("--- %T\n", v.Callee)
				funcType := computeType(v.Callee, a)
				ftyp, ok := funcType.(*nodes.FuncType)
				if !ok {
					errs = append(errs, fmt.Sprintf("Type is not callable: %s", funcType.Text()))
					return errs
				}
				// if !ftyp.ReturnType.Equals(&nodes.VoidType{}) {
				// 	errs = append(errs, fmt.Sprintf("Return value of type %s from function %s is not used", ftyp.ReturnType.Text(), "<TODO>"))
				// 	return errs
				// }
				// all args are of expected type
				if len(v.Args) != len(ftyp.ArgTypes) {
					errs = append(errs, fmt.Sprintf("Function %s expects %d parameters, but supplied %d", v.Callee.String(""), len(ftyp.ArgTypes), len(v.Args)))
					return errs
				}
				for i, k := range v.Args {
					expT := ftyp.ArgTypes[i]
					passedT := computeType(k, a)
					// fmt.Println(expT.Text(), passedT.Text())
					if !expT.Equals(passedT) {
						errs = append(errs, fmt.Sprintf("%d th parameter to function %s should be %s, not %s", i, "<TODO>", expT.Text(), passedT.Text()))
					}
				}

			}
		default:
			errs = append(errs, fmt.Sprintf("Can't check for %T", v))
		}
	}
	return errs
}
