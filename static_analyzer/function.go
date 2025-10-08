package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	"he++/utils"
)

func (a *Analyzer) registerFunctionDecl(fnd *nodes.FuncNode) {
	// todo: if supporting function overloading
	// key should have args types too
	a.definedSyms[fnd.Name] = computeType(fnd, a)
}

func (a *Analyzer) checkFunctionDef(fnd *nodes.FuncNode) {
	for _, arg := range fnd.ArgList {
		a.definedSyms[arg.Name] = arg.DataT
	}
	// todo: instead of passing returnType, look up the scope stack
	// to see what function we're inside. (todo: scope stack)
	a.checkScope(fnd.Scope, fnd.ReturnType)
	stmts := fnd.Scope.Children
	returnFound := false
	
	for _, stmt := range stmts {
		if stmt.Type() == nodes.RETURN {
			returnFound = true
		}
	}
	if !returnFound {
		if _, ok := fnd.ReturnType.(*nodes.VoidType); !ok {
			a.AddError(
				fnd.Scope.Range().End,
				utils.TypeError,
				fmt.Sprintf("Expected return value of type %s", fnd.ReturnType.Text()),
			)
		}
	}
	return
}
