package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	"he++/utils"
	// "he++/utils"
)

var _ = fmt.Println

/**
Things to check for:
- idents not used before decl, no duplicate decl
- no expressions (OPERATOR) in top level scope
- no break, continue in non loop scope
- all code paths should have return
- type inference in var decl
- type consistency in func return
- normalized types
*/

type Analyzer struct {
	scopeStack utils.Stack[ScopeEntry]
	// refers to data types
	// todo: would it be OK to merge with syms?
	definedTypes map[string]*utils.Stack[nodes.DataType]
	// sym refers to functions and variables
	definedSyms map[string]*utils.Stack[nodes.DataType]
	Errs        []utils.CompilerError
}

func MakeAnalyzer() Analyzer {
	a := Analyzer{
		scopeStack:   utils.MakeStack[ScopeEntry](),
		definedTypes: make(map[string]*utils.Stack[nodes.DataType]),
		definedSyms:  make(map[string]*utils.Stack[nodes.DataType]),
	}
	a.PushScope(BASE)
	addInbuiltDefinitions(&a)
	return a
}

func (a *Analyzer) PushScope(st ScopeType) {
	scopeEntry := ScopeEntry{ScopeType: st, DefinedSyms: make(map[string]bool), DefinedTypes: make(map[string]bool)}
	a.scopeStack.Push(scopeEntry)
}

func (a *Analyzer) PopScope() {
	lastScope := a.GetLatestScope()
	for k := range lastScope.DefinedSyms {
		v := a.definedSyms[k]
		v.Pop()
		a.definedSyms[k] = v
	}
	for k := range lastScope.DefinedTypes {
		k := a.definedTypes[k]
		k.Pop()
	}
	a.scopeStack.Pop()
}

func (a *Analyzer) DefineType(name string, dt nodes.DataType) {
	lastScope := a.GetLatestScope()
	lastScope.DefinedTypes[name] = true
	_, exists := a.definedTypes[name]
	if !exists {
		stk := utils.MakeStack[nodes.DataType]()
		a.definedTypes[name] = &stk
	}
	a.definedTypes[name].Push(dt)
}

func (a *Analyzer) DefineSym(name string, dt nodes.DataType) {
	lastScope := a.GetLatestScope()
	lastScope.DefinedSyms[name] = true
	_, exists := a.definedSyms[name]
	if !exists {
		stk := utils.MakeStack[nodes.DataType]()
		a.definedSyms[name] = &stk
	}
	a.definedSyms[name].Push(dt)
}

func (a *Analyzer) GetSym(key string) (nodes.DataType, bool) {
	symStk, exists := a.definedSyms[key]
	if !exists {
		return &nodes.ErrorType{}, false
	}
	val, exists := symStk.Peek()
	if !exists {
		return &nodes.ErrorType{}, false
	}
	return val, true
}

func (a *Analyzer) GetLatestScope() ScopeEntry {
	lastScope, exists := a.scopeStack.Peek()
	if !exists {
		panic("Stack shouldn't have been empty here!")
	}
	return lastScope
}

func (a *Analyzer) AnalyzeAST(n *nodes.SourceFileNode) {

	for _, ch := range n.Children {
		if ch.Type() == nodes.FUNCTION {
			funcNode := ch.(*nodes.FuncNode)
			a.registerFunctionDecl(funcNode)
		}
	}

	for _, ch := range n.Children {
		if funcNode, ok := ch.(*nodes.FuncNode); ok {
			a.checkFunctionDef(funcNode)
		}
	}
	fmt.Printf("In source file %s:\n", utils.Underline(n.FilePath))
	for _, k := range a.Errs {
		fmt.Println(&k)
	}
	return
}

func (a *Analyzer) AddError(Line int, Name utils.CompilerErrorKind, Msg string) {
	a.Errs = append(a.Errs, utils.CompilerError{Line: Line, Name: Name, Msg: Msg})
}
