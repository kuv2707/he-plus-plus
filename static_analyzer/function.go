package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
)

func (a *Analyzer) registerFunctionDecl(fnd *nodes.FuncNode) {
	// todo: if supporting function overloading
	// key should have args types too
	a.definedSyms[fnd.Name] = computeType(fnd, a)
}

func (a *Analyzer) checkFunctionDef(fnd *nodes.FuncNode) []string {
	for _, arg := range fnd.ArgList {
		a.definedSyms[arg.Name] = arg.DataT
	}
	errs := a.checkScope(fnd.Scope)
	stmts := fnd.Scope.Children
	returnChecked := false
	// if _, ok := fnd.ReturnType.(*nodes.NamedType)

	// if _, exists := a.definedTypes[fnd.ReturnType.Text()]; !exists {
	// 	errs = append(errs, fmt.Sprint("Type not defined: ", fnd.ReturnType))
	// }
	
	for _, stmt := range stmts {
		if stmt.Type() == nodes.RETURN {
			ret := stmt.(*nodes.ReturnNode)
			typ := computeType(ret.Value, a)
			// todo: add func in nodes.DataType to compare two types
			if typ.Text() != fnd.ReturnType.Text() {
				errs = append(errs, fmt.Sprintf("Expected to return %s, found %s", fnd.ReturnType.Text(), typ.Text()))
			}
			returnChecked = true
		}
	}
	if !returnChecked {
		if _, ok := fnd.ReturnType.(*nodes.VoidType); !ok {
			errs = append(errs, fmt.Sprintf("Expected return value of type %s", fnd.ReturnType.Text()))
		}
	}
	return errs
}
