package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	"he++/utils"
)

func (a *Analyzer) checkExpression(exp nodes.TreeNode) nodes.DataType {
	switch v := exp.(type) {
	case *nodes.NumberNode:
		switch v.NumType {
		case nodes.FLOAT_NUMBER:
			return &FLOAT_DATATYPE
		case nodes.INT32_NUMBER:
			return &INT_DATATYPE
		}
	case *nodes.InfixOperatorNode:
		l := a.checkExpression(v.Left)
		r := a.checkExpression(v.Right)
		// todo: check if l and r are compatible under this optype
		return a.operatorReturnType(v.Op, l, r)
	case *nodes.IdentifierNode:
		varname := v.Name()
		s, exists := a.GetSym(varname)
		if !exists {
			a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Undefined identifier %s in expression", utils.Green(varname)))
			return &ERROR_TYPE
		}
		return s
	default:
		a.AddError(
			v.Range().Start,
			utils.UndefinedError,
			fmt.Sprintf("Can't check for expresion node %T", v),
		)
	}
	return &ERROR_TYPE
}
