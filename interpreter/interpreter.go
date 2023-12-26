package interpreter

import (
	"fmt"
	"toylingo/parser"
)

func Interpret(root *parser.TreeNode) {
	pushScopeContext("scope","root")
	ctx:=contextStack.GetStack()[0].(scopeContext)
	addNativeFuncDeclarations(&ctx)
	// for _,v := range ctx.functions{
	// 	v.PrintTree("->")
	// }
	executeScope(root, &ctx)
	printMemoryStats()
}



const REASON_NATURAL Reason = "natural"
const REASON_RETURN Reason = "return"
const REASON_BREAK Reason = "break"

const TYPE_LOOP string = "loop"
const TYPE_FUNCTION string = "function"
const TYPE_SCOPE string = "scope"
const TYPE_CONDITIONAL string = "conditional"

func executeScope(node *parser.TreeNode, ctx *scopeContext) (Reason, *Variable) {
	debug_info("entered", ctx.scopeName)
	var returnReason Reason = REASON_NATURAL
	scopeType := ctx.scopeTyp
SCOPE_EXECUTION:
	for i := range node.Children {
		child := node.Children[i]
		LineNo = child.LineNo
		switch child.Label {
		case "function":
			ctx.functions[child.Description] = *child

		case "scope":
			rzn, val := executeScope(child, pushScopeContext(TYPE_SCOPE,"simple_scope"))
			ctx.returnValue = val
			if rzn != REASON_NATURAL {
				if rzn == REASON_BREAK {
					break SCOPE_EXECUTION
				} else if scopeType == TYPE_FUNCTION && rzn == REASON_RETURN {

					break SCOPE_EXECUTION
				} else {
					returnReason = rzn
					break SCOPE_EXECUTION
				}
			}
		case "conditional_block":
			debug_info("conditional block")
			k := 0
			executed:=false
			for ; ; k++ {
				condnode, exists := child.Properties["condition"+fmt.Sprint(k)]
				if !exists {
					break
				}
				condition := evaluateExpressionClean(condnode, ctx)
				if condition == 0 {
					continue
				}
				rzn, val := executeScope(child.Properties["ifnode"+fmt.Sprint(k)], pushScopeContext(TYPE_CONDITIONAL,"if-elif"))
				ctx.returnValue = val
				if rzn != REASON_NATURAL {
					if rzn == REASON_BREAK {
						returnReason = rzn
						break SCOPE_EXECUTION
					} else if scopeType == TYPE_FUNCTION && rzn == REASON_RETURN {
						break SCOPE_EXECUTION
					} else {
						returnReason = rzn
						break SCOPE_EXECUTION
					}
				}
				executed=true
				break
			}
			if child.Properties["else"] == nil || executed{
				continue SCOPE_EXECUTION
			}
			rzn, val := executeScope(child.Properties["else"], pushScopeContext(TYPE_CONDITIONAL,"else"))
			ctx.returnValue = val
			if rzn != REASON_NATURAL {
				if rzn == REASON_BREAK {
					returnReason = rzn
					break SCOPE_EXECUTION
				} else if scopeType == TYPE_FUNCTION && rzn == REASON_RETURN {
					break SCOPE_EXECUTION
				} else {
					returnReason = rzn
					break SCOPE_EXECUTION
				}
			}

		case "loop":
			for true {
				num := evaluateExpressionClean(child.Properties["condition"], ctx)
				if num == 0 {
					break
				}
				rzn, val := executeScope(child.Properties["body"], pushScopeContext(TYPE_LOOP,"lp"))
				if rzn == REASON_BREAK {
					break
				}
				ctx.returnValue = val
				if rzn != REASON_NATURAL {
					if rzn == REASON_BREAK {
						//todo: match loop label
						returnReason = rzn
						break SCOPE_EXECUTION
					} else if scopeType == TYPE_FUNCTION && rzn == REASON_RETURN {
						break SCOPE_EXECUTION
					} else {
						returnReason = rzn
						break SCOPE_EXECUTION
					}

				}
			}
		case "operator":
			evaluateExpressionClean(child, ctx)
		case "call":

			evaluateFuncCall(*child, ctx)
			gc()
		case "break":
			returnReason = REASON_BREAK
			break SCOPE_EXECUTION

		case "return":
			expr := evaluateExpression(child.Children[0], ctx)
			expr.pointer.temp = false
			ctx.returnValue = &expr
			returnReason = REASON_RETURN
			break SCOPE_EXECUTION
		default:
			debug_error("__unknown", child.Label)
		}

	}
	printMemoryStats()
	debug_info("exited", ctx.scopeName)
	popScopeContext()
	return returnReason, ctx.returnValue
}

// will only return number value from evaluated variable
func evaluateExpressionClean(node *parser.TreeNode, ctx *scopeContext) float64 {
	variable := evaluateExpression(node, ctx)
	ret := getValue(variable)
	gc()
	return ret
}
