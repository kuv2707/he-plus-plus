package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	"he++/utils"
)

// partners with computeType
func (a *Analyzer) checkExpression(exp nodes.TreeNode) nodes.DataType {
	exp.String(&utils.ASTPrinter{})
	switch v := exp.(type) {
	case *nodes.NumberNode:
		switch v.NumType {
		case nodes.FLOAT_NUM:
			return FLOAT_DATATYPE
		case nodes.INT_NUM:
			return INT_DATATYPE
		default:
			return ERROR_TYPE
		}
	case *nodes.BooleanNode:
		return BOOLEAN_DATATYPE
	case *nodes.InfixOperatorNode:
		l := a.computeType(v.Left)
		r := a.computeType(v.Right)
		// todo: check if l and r are compatible under this optype
		ort := a.operatorReturnType(v.Op, l, r, v.Range().Start)
		if isErrorType(ort) {
			a.AddError(v.Range().Start, utils.TypeError, fmt.Sprintf("Can't perform %s on types %s and %s", v.Op, utils.Cyan(l.Text()), utils.Cyan(r.Text())))
		}
		v.ResultDT = ort
		return ort
	case *nodes.IdentifierNode:
		varname := v.Name()
		s, exists, readAs := a.GetSymInfo(varname)
		v.ChangeName(readAs)
		if !exists {
			a.AddError(v.Range().Start, utils.UndefinedError, fmt.Sprintf("Undefined identifier %s in expression", utils.Green(varname)))
			return ERROR_TYPE
		}
		s.numUses += 1
		v.DataT = s.dt
		return s.dt
	case *nodes.FuncCallNode:
		a.computeType(v)
	default:
		a.AddError(
			v.Range().Start,
			utils.UndefinedError,
			fmt.Sprintf("Can't check for expresion node %T", v),
		)
	}
	return ERROR_TYPE
}
