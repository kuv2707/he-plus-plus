package staticanalyzer

import (
	"fmt"
	nodes "he++/parser/node_types"
	"he++/utils"
	// "he++/utils"
)

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

type VarDefInfo struct {
	dt      nodes.DataType
	numUses int
}

type Analyzer struct {
	scopeStack utils.Stack[ScopeEntry]
	// refers to data types
	// todo: would it be OK to merge with syms?
	definedTypes map[string]*nodes.DataType
	// sym refers to functions and variables
	definedSyms map[string]*VarDefInfo
	// operatorTypeRelations map[nodes.TypeId]
	Errs []utils.CompilerError
}

func MakeAnalyzer() Analyzer {
	a := Analyzer{
		scopeStack:   *utils.MakeStack[ScopeEntry](),
		definedTypes: make(map[string]*nodes.DataType),
		definedSyms:  make(map[string]*VarDefInfo),
	}
	a.PushScope(BASE)
	addFundamentalDefinitions(&a)
	return a
}

func (a *Analyzer) PushScope(st ScopeType) {
	scopeEntry := ScopeEntry{ScopeType: st, DefinedSyms: make(map[string]bool), DefinedTypes: make(map[string]bool), symRedirects: make(map[string]string)}
	a.scopeStack.Push(scopeEntry)
}

func (a *Analyzer) PopScope() {
	lastScope := a.GetLatestScope()
	for k := range lastScope.DefinedSyms {
		delete(a.definedSyms, k)
	}
	for k := range lastScope.DefinedTypes {
		delete(a.definedTypes, k)
	}
	a.scopeStack.Pop()
}

// todo: uncomment when supporting struct defns in validation
// func (a *Analyzer) DefineType(name string, dt nodes.DataType) {
// 	lastScope := a.GetLatestScope()
// 	lastScope.DefinedTypes[name] = true
// 	_, exists := a.definedTypes[name]
// 	if !exists {
// 		a.definedTypes[name] = utils.MakeStack[nodes.DataType]()
// 	}
// 	a.definedTypes[name].Push(dt)
// }

func (a *Analyzer) GetType(key string) (nodes.DataType, bool) {
	val, exists := a.definedTypes[key]
	if !exists {
		return ERROR_TYPE, false
	}

	return *val, true
}

func (a *Analyzer) DefineSym(name string, dt nodes.DataType) string {
	// maybe instead of stack, just change the name of this var to sth unique
	// and redirect all references in the scope to the new name

	lastScope := a.GetLatestScope()
	_, exists := a.definedSyms[name]
	if exists {
		// store it as a different name to avoid collision with previous definitions in outer scopes
		// since it begins with a number, it won't collide with user defined vars
		newName := fmt.Sprintf("%d%s", a.scopeStack.Len(), name)
		lastScope.symRedirects[name] = newName
		name = newName
	}
	lastScope.DefinedSyms[name] = true
	a.definedSyms[name] = &VarDefInfo{dt, 0}
	return name
}

func (a *Analyzer) afterRedirect(key string) string {
	lastScope := a.GetLatestScope()
	if k, ex := lastScope.symRedirects[key]; ex {
		return k
	}
	return key
}

func (a *Analyzer) GetSymInfo(key string) (*VarDefInfo, bool, string) {
	key = a.afterRedirect(key)
	symInfo, exists := a.definedSyms[key]
	if !exists {
		return nil, false, key
	}
	return symInfo, true, key
}

func (a *Analyzer) GetLatestScope() *ScopeEntry {
	lastScope, exists := a.scopeStack.Peek()
	if !exists {
		panic("Stack shouldn't have been empty here!")
	}
	return lastScope
}

func (a *Analyzer) AnalyzeAST(n *nodes.SourceFileNode) bool {

	for _, ch := range n.Children {
		if ch.Type() == nodes.FUNCTION {
			funcNode := ch.(*nodes.FuncNode)
			a.registerFunctionDecl(funcNode)
		}
	}

	// todo: use parallel iterator
	for _, ch := range n.Children {
		if funcNode, ok := ch.(*nodes.FuncNode); ok {
			a.checkFunctionDef(funcNode)
		}
	}
	fmt.Printf("In source file %s:\n", utils.Underline(n.FilePath))
	for _, k := range a.Errs {
		fmt.Println(&k)
	}
	return len(a.Errs) == 0
}

func (a *Analyzer) AddError(Line int, Name utils.CompilerErrorKind, Msg string) {
	a.Errs = append(a.Errs, utils.CompilerError{Line: Line, Name: Name, Msg: Msg})
}
