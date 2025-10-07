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
	a.checkScope(fnd.Scope)
	stmts := fnd.Scope.Children
	returnChecked := false

	for _, stmt := range stmts {
		if stmt.Type() == nodes.RETURN {
			ret := stmt.(*nodes.ReturnNode)
			typ := computeType(ret.Value, a)
			// todo: add func in nodes.DataType to compare two types
			if !typ.Equals(fnd.ReturnType) {
				a.AddError(
					ret.Range().Start,
					utils.TypeError,
					fmt.Sprintf("Expected to return %s, found %s", fnd.ReturnType.Text(), typ.Text()),
				)
			}
			returnChecked = true
		}
	}
	if !returnChecked {
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
