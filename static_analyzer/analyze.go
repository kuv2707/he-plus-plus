package staticanalyzer

import (
	nodes "he++/parser/node_types"
	"he++/utils"
)

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

func (a *Analyzer) analyzeAST(n *nodes.TreeNode) {

}
