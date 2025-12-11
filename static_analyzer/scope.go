package staticanalyzer

import (
	"fmt"
	"he++/lexer"
	nodes "he++/parser/node_types"
	"he++/utils"
)

func (a *Analyzer) checkScope(scp *nodes.ScopeNode) nodes.DataType {
	var scopeRet nodes.DataType = nil
	// check if datatypes exist and var decls match type
	for i, n := range scp.Children {
		scopeRet = a.checkNode(n, i, scp, scopeRet)
	}
	// todo: before popping, the following info collected while analyzing
	// the scope should be persisted in the node:
	// total num of vars defined (used to calculate frame size)
	// number of usages of each identifier (for register allotment)
	// ...
	a.PopScope()

	return scopeRet
}

func (a *Analyzer) checkNode(n nodes.TreeNode, i int, scp *nodes.ScopeNode, scopeRet nodes.DataType) nodes.DataType {
	switch v := n.(type) {
	case *nodes.VariableDeclarationNode:
		{
			for _, tn := range v.Declarations {
				if op := tn.(*nodes.InfixOperatorNode); op.Op == lexer.ASSN {
					varname := op.Left.(*nodes.IdentifierNode)
					// todo: Check if v.DataT itself is valid
					// rval should have same type
					rvalType := a.computeType(op.Right)
					a.DefineSym(varname.Name(), v.DataT)
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
						fmt.Sprintf("%s not allowed. Use %s", utils.Red(op.Op), utils.Green(lexer.ASSN)),
					)
				}
			}
		}
	case *nodes.ReturnNode:
		{
			scopeRet = a.computeType(v.Value)
			if i < len(scp.Children)-1 {
				// no esperamos que haya mas nudos a procesar
				a.AddError(v.Range().End, utils.SyntaxError, "A return statement must be the last statement in the scope.")
			}

		}
	case *nodes.FuncCallNode:
		{
			// todo: maybe show a warning if the type returned isn't void, meaning return value never used
			a.computeType(v)
		}
	case *nodes.ScopeNode:
		{
			a.PushScope(NESTED)
			ret := a.checkScope(v)
			if scopeRet != nil && ret != nil {
				if !scopeRet.Equals(ret) {
					// the scope had earlier returned `scopeRet`, but now seems to return `ret`
					// lets take the earlier return type to be the expected one
					a.AddError(v.Range().End, utils.TypeError, fmt.Sprintf("Expected return value of type %s, got %s", scopeRet.Text(), ret.Text()))
				}
			} else if scopeRet != nil && ret == nil {
				// esta bien, significa que este scope en particular no devuelve nada
				// asi que no hay ningun modo de causar un error
			} else if scopeRet == nil && ret != nil {
				// digamos que este es el tipo de retorno
				scopeRet = ret
			}
			// el cuarto caso es como el segundo
		}
	case *nodes.IfNode:
		{
			exhaustive := true
			var retType nodes.DataType
			for _, branch := range v.Branches {
				condTyp := a.computeType(branch.Condition)
				if !isBooleanType(condTyp) {
					a.AddError(branch.Condition.Range().Start, utils.TypeError,
						fmt.Sprintf("Expected the expression to evaluate to %s or %s", utils.Blue(lexer.TRUE), utils.Blue(lexer.FALSE)))
				}
				a.PushScope(CONDITIONAL)
				ret := a.checkScope(branch.Scope)
				if ret == nil {
					exhaustive = false
				} else {

					if retType == nil {
						retType = ret
					} else {
						if retType != ret {
							// the return types of the scopes don't agree
							a.AddError(branch.Scope.Range().Start, utils.TypeError,
								fmt.Sprintf("Expected to return %s or nothing", utils.Cyan(retType.Text())))
						}
					}

				}
			}
			v.Exhaustive = exhaustive
		}
	case *nodes.LoopNode:
		{
			a.PushScope(LOOP)
			scopeRet = a.checkNode(v.Initializer, i, scp, scopeRet)
			condTyp := a.computeType(v.Condition)
			if !isBooleanType(condTyp) {
				a.AddError(v.Condition.Range().Start, utils.TypeError,
					fmt.Sprintf("Expected the expression to evaluate to %s or %s", utils.Blue(lexer.TRUE), utils.Blue(lexer.FALSE)))
			}
			a.checkNode(v.Updater, i, scp, scopeRet)
			scopeRet = a.checkScope(v.Scope) // todo: propagate return value
		}
	case *nodes.InfixOperatorNode:
		if v.Op != lexer.ASSN {
			a.AddError(v.Range().Start, utils.NotAllowed, fmt.Sprintf("Unused expression result for op %s", v.Op))
		} else {
			a.checkExpression(v)
		}
	case *nodes.EmptyPlaceholderNode:
		{
			// no hacer nada
		}
	default:
		a.AddError(
			v.Range().Start,
			utils.UndefinedError,
			fmt.Sprintf("Can't perform static analysis for %T", v),
		)
	}
	return scopeRet
}
