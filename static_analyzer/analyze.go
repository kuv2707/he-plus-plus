package staticanalyzer

import "he++/utils"

/**
Things to check for:
- idents not used before decl, no duplicate decl
- no expressions (OPERATOR) in top level scope
- no break, continue in non loop scope
- all code paths should have return
- type inference in var decl
- type consistency in func return
*/

type Analyzer struct {
	scopeStack utils.Stack
}

func (a *Analyzer) PushScope() {
	a.scopeStack.Push(0)
}

func MakeAnalyzer() Analyzer {
	return Analyzer{utils.MakeStack()}
}
