package interpreter

import (
	"fmt"
	// "he++/globals"
	"he++/parser"
	"he++/utils"
)

func Init(a map[string]string) *ScopeContext {
	args = a
	pointers[0] = NULL_POINTER
	ctx := pushScopeContext("scope", "root")
	addNativeFuncDeclarations(ctx)
	return ctx
}

var args map[string]string

func Interpret(root *parser.TreeNode, ctx *ScopeContext) *ScopeContext {
	executeScope(root, ctx)
	printMemoryStats()
	return ctx
}

const REASON_NATURAL Reason = "natural"
const REASON_RETURN Reason = "return"
const REASON_BREAK Reason = "break"

const TYPE_LOOP string = "loop"
const TYPE_FUNCTION string = "function"
const TYPE_SCOPE string = "scope"
const TYPE_CONDITIONAL string = "conditional"

func executeScope(node *parser.TreeNode, ctx *ScopeContext) (Reason, *Pointer) {
	debug_info("entered", ctx.scopeId)
	// printMemoryStats()
	var returnReason Reason = REASON_NATURAL
	scopeType := ctx.scopeType
SCOPE_EXECUTION:
	for i := range node.Children {
		child := node.Children[i]
		ctx.currentLine = child.LineNo
		switch child.Label {
		case "function":
			ctx.functions[child.Description] = *child

		case "scope":
			rzn, val := executeScope(child, pushScopeContext(TYPE_SCOPE, "simple_scope"))
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
			//debug_info("conditional block")
			k := 0
			executed := false
			for ; ; k++ {
				condnode, exists := child.Properties["condition"+fmt.Sprint(k)]
				if !exists {
					break
				}
				res := evaluateExpression(condnode, ctx)
				result := booleanValue(res)
				freePtr(res)
				if !result {
					continue
				}
				rzn, val := executeScope(child.Properties["ifnode"+fmt.Sprint(k)], pushScopeContext(TYPE_CONDITIONAL, "if-elif"))
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
				executed = true
				break
			}
			if child.Properties["else"] == nil || executed {
				continue SCOPE_EXECUTION
			}
			rzn, val := executeScope(child.Properties["else"], pushScopeContext(TYPE_CONDITIONAL, "else"))
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
			loopctx := makeScopeContext(TYPE_LOOP, "loop")
			for true {
				res := evaluateExpression(child.Properties["condition"], ctx)
				result := booleanValue(res)
				if !result {
					break
				}
				for k := range loopctx.functions {
					delete(loopctx.functions, k)
				}
				for k := range loopctx.variables {
					delete(loopctx.variables, k)
				}
				pushToContextStack(loopctx)
				rzn, val := executeScope(child.Properties["body"], &loopctx)
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
			if !utils.IsOneOf(child.Description, []string{"++", "--", "="}) {
				interrupt(child.LineNo, "Operator", child.Description, "is not allowed in statements")
			}
			evaluateExpression(child, ctx)
		case "call":

			garb := evaluateFuncCall(*child, ctx)
			if !garb.isNull() {
				freePtr(garb)
			}
		case "break":
			returnReason = REASON_BREAK
			break SCOPE_EXECUTION

		case "return":
			expr := evaluateExpression(child.Children[0], ctx)
			expr.temp = false
			ctx.returnValue = expr
			//debug_info("return value is", expr)
			returnReason = REASON_RETURN
			break SCOPE_EXECUTION
		default:
			debug_error("__unknown", child.Label)
		}
	}
	gc()
	// printMemoryStats()
	debug_info("exited", ctx.scopeId)
	popScopeContext()
	return returnReason, ctx.returnValue
}
