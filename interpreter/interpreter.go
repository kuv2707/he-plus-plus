package interpreter

import (
	"fmt"
	"strings"
	"toylingo/parser"
)

func Interpret(root *parser.TreeNode) {
	executeScope(root, pushScopeContext("scope_0"))
	printMemoryStats()
}

type Reason string

const REASON_NATURAL Reason = "natural"
const REASON_RETURN Reason = "return"
const REASON_BREAK Reason = "break"

const TYPE_LOOP string = "loop"
const TYPE_FUNCTION string = "function"
const TYPE_SCOPE string = "scope"

func executeScope(node *parser.TreeNode, ctx *scopeContext) (Reason, *Variable) {
	fmt.Println("entered", ctx.scopeType)
	var returnReason Reason = REASON_NATURAL
	scopeType := strings.Split(ctx.scopeType, "_")[0]
SCOPE_EXECUTION:
	for i := range node.Children {
		child := node.Children[i]
		switch child.Label {
		case "function":
			ctx.functions[child.Description] = *child

		case "scope":
			rzn, val := executeScope(child, pushScopeContext(TYPE_SCOPE))
			ctx.returnValue = val
			fmt.Println("scope ret rzn", rzn, scopeType)
			if rzn != REASON_NATURAL {
				if scopeType == TYPE_LOOP && rzn == REASON_BREAK {
					break SCOPE_EXECUTION
				} else if scopeType == TYPE_FUNCTION && rzn == REASON_RETURN {

					break SCOPE_EXECUTION
				} else {
					returnReason = rzn
					break SCOPE_EXECUTION
				}
			}
		case "conditional_block":
			k := 0
			for ; ; k++ {
				condnode,exists:=child.Properties["condition"+fmt.Sprint(k)]
				if !exists{
					break
				}
				condition := evaluateExpressionClean(condnode, ctx)
				if condition == 0 {
					continue
				}
				rzn, val := executeScope(child.Properties["ifnode"+fmt.Sprint(k)], pushScopeContext(TYPE_SCOPE))
				ctx.returnValue = val
				if rzn != REASON_NATURAL {
					if scopeType == TYPE_LOOP && rzn == REASON_BREAK {
						break SCOPE_EXECUTION
					} else if scopeType == TYPE_FUNCTION && rzn == REASON_RETURN {
						break SCOPE_EXECUTION
					} else {
						returnReason = rzn
						break SCOPE_EXECUTION
					}
				}
				break
			}
		case "loop":
			for true {
				num := evaluateExpressionClean(child.Properties["condition"], ctx)
				if num == 0 {
					break
				}
				rzn, val := executeScope(child.Properties["body"], pushScopeContext(TYPE_LOOP))
				ctx.returnValue = val
				if rzn != REASON_NATURAL {
					if scopeType == TYPE_LOOP && rzn == REASON_BREAK {
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
			fallthrough
		case "call":
			evaluateExpressionClean(child, ctx)

		case "return":
			expr := evaluateExpression(child.Children[0], ctx)
			expr.pointer.temp = false
			ctx.returnValue = &expr
			returnReason = REASON_RETURN
			break SCOPE_EXECUTION
		default:
			fmt.Println("__unknown", child.Label)
		}

	}
	printMemoryStats()
	popScopeContext()
	return returnReason, ctx.returnValue
}

// will only return number value from evaluated variable
func evaluateExpressionClean(node *parser.TreeNode, ctx *scopeContext) float64 {
	variable := evaluateExpression(node, ctx)
	ret := 0.0
	if variable.vartype == TYPE_NUMBER || variable.vartype == TYPE_BOOLEAN {
		ret = getNumber(variable)
	}
	gc()
	return ret
}
