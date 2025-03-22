package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
)

func (a *Analyzer) registerFunctionDecl(fnd *nodes.FuncNode) {
	// todo: if supporting function overloading
	// key should have args types too
	a.functionDecls[fnd.Name] = fnd
}

func (a *Analyzer) checkFunctionDef(fnd *nodes.FuncNode) []string {
	errs := a.checkScope(fnd.Scope)
	stmts := fnd.Scope.Children
	returnChecked := false
	if _, exists := a.definedTypes[fnd.ReturnType.Text]; !exists {
		errs = append(errs, fmt.Sprint("Type not defined: ", fnd.ReturnType))
	}
	for _, stmt := range stmts {
		if stmt.Type() == nodes.RETURN {
			ret := stmt.(*nodes.ReturnNode)
			typ := computeType(ret.Value, a)
			// todo: add func in nodes.DataType to compare two types
			if typ.Text != fnd.ReturnType.Text {
				errs = append(errs, fmt.Sprintf("Expected to return %s, found %s", fnd.ReturnType, typ.Text))
			}
			returnChecked = true
		}
	}
	if !returnChecked {
		if fnd.ReturnType.Text != "void" {
			errs = append(errs, fmt.Sprintf("Expected return value of type %s", fnd.ReturnType))
		}
	}
	return errs
}
