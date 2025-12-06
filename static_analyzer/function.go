package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	"he++/utils"
)

func (a *Analyzer) registerFunctionDecl(fnd *nodes.FuncNode) {
	// todo: if supporting function overloading,
	// then the key should have args types too
	a.DefineSym(fnd.Name, a.computeType(fnd))
}

func (a *Analyzer) checkFunctionDef(fnd *nodes.FuncNode) {
	a.PushScope(FUNCTION)
	for _, arg := range fnd.ArgList {
		a.DefineSym(arg.Name, arg.DataT)
	}
	// todo: instead of passing returnType, look up the scope stack
	// to see what function we're inside. (todo: scope stack)
	ret := a.checkScope(fnd.Scope)
	if ret == nil {
		ret = &nodes.VOID_DATATYPE
	}
	if !fnd.ReturnType.Equals(ret) {
		a.AddError(fnd.Range().Start, utils.TypeError,
			fmt.Sprintf("Expected to return value of type %s but found %s", utils.Cyan(fnd.ReturnType.Text()), utils.Cyan(ret.Text())))
	}

	return
}
