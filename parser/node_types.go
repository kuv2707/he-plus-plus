package parser

import (
	"fmt"
	"he++/globals"
)

type TreeNodeType string

const (
	SCOPE       TreeNodeType = "Scope"
	CONDITIONAL TreeNodeType = "Conditional"
	LOOP        TreeNodeType = "Loop"
	FUNCTION    TreeNodeType = "Function"
	STRUCT      TreeNodeType = "Struct"
	OPERATOR    TreeNodeType = "Expression"
	VALUE       TreeNodeType = "Value"
	ARR_IND     TreeNodeType = "Array_Index"
	VAR_DECL    TreeNodeType = "Variable_Declaration"
	RETURN      TreeNodeType = "Return"
)

const TAB = "  "

type TreeNode interface {
	String(ind string) string
	Type() TreeNodeType
}

type EmptyPlaceholderNode struct {
}

func (e *EmptyPlaceholderNode) String(ind string) string {
	return ind + "<empty>"
}

func (e *EmptyPlaceholderNode) Type() TreeNodeType {
	return TreeNodeType("")
}

// expression related nodes
type NumberNode struct {
	dataBytes []byte
	numType   string
}

func (n *NumberNode) String(ind string) string {
	return ind + string(n.dataBytes)
}

func (n *NumberNode) Type() TreeNodeType {
	return VALUE
}

func NewNumberNode(dataBytes []byte, numType string) *NumberNode {
	return &NumberNode{dataBytes, numType}
}

type StringNode struct {
	dataBytes []byte
}

func (s *StringNode) String(ind string) string {
	return ind + string(s.dataBytes)
}

func (s *StringNode) Type() TreeNodeType {
	return VALUE
}

func NewStringNode(dataBytes []byte) *StringNode {
	return &StringNode{dataBytes}
}

type BooleanNode struct {
	dataBytes []byte
}

func (b *BooleanNode) String(ind string) string {
	return ind + string(b.dataBytes)
}

func (b *BooleanNode) Type() TreeNodeType {
	return VALUE
}

func NewBooleanNode(dataBytes []byte) *BooleanNode {
	return &BooleanNode{dataBytes}
}

type IdentifierNode struct {
	name string
}

func (i *IdentifierNode) Name() string {
	return i.name
}

func (i *IdentifierNode) String(ind string) string {
	return ind + i.name
}

func (i *IdentifierNode) Type() TreeNodeType {
	return VALUE
}

func NewIdentifierNode(name string) *IdentifierNode {
	return &IdentifierNode{name}
}

const (
	PREFIX  = "pre"
	POSTFIX = "post"
)

type PrePostOperatorNode struct {
	opType  string
	op      string
	operand TreeNode
}

func (o *PrePostOperatorNode) String(ind string) string {
	return fmt.Sprintf("%s%s %s\n%s", ind, o.opType, o.op, o.operand.String(ind+TAB))
}

func (o *PrePostOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewPrePostOperatorNode(opType string, op string, operand TreeNode) *PrePostOperatorNode {
	return &PrePostOperatorNode{opType, op, operand}
}

type InfixOperatorNode struct {
	left  TreeNode
	op    string
	right TreeNode
}

func (o *InfixOperatorNode) String(ind string) string {
	return ind + o.op + "\n" + o.left.String(ind+TAB) + "\n" + o.right.String(ind+TAB)
}

func (o *InfixOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewInfixOperatorNode(left TreeNode, op string, right TreeNode) *InfixOperatorNode {
	return &InfixOperatorNode{left, op, right}
}

type TernaryOperatorNode struct {
	condition TreeNode
	ifTrue    TreeNode
	ifFalse   TreeNode
}

func (t *TernaryOperatorNode) String(ind string) string {
	return ind + "ternary\n" + t.condition.String(ind+TAB) + "\n" + t.ifTrue.String(ind+TAB) + "\n" + t.ifFalse.String(ind+TAB)
}

func (t *TernaryOperatorNode) Type() TreeNodeType {
	return OPERATOR
}

func NewTernaryNode(condition TreeNode, ifTrue TreeNode, ifFalse TreeNode) *TernaryOperatorNode {
	return &TernaryOperatorNode{condition, ifTrue, ifFalse}
}

type FuncCallNode struct {
	callee TreeNode
	args   []TreeNode
}

func (f *FuncCallNode) String(ind string) string {
	args := ""
	for _, arg := range f.args {
		args += arg.String(ind+TAB+"  ") + "\n"
	}
	return ind + "call\n" + f.callee.String(ind+TAB) + ind + "\n" + ind + "  args:\n" + args
}

func (f *FuncCallNode) Type() TreeNodeType {
	return FUNCTION
}

func NewFuncCallNode(name TreeNode) *FuncCallNode {
	return &FuncCallNode{name, make([]TreeNode, 0)}
}

func (f *FuncCallNode) arg(arg TreeNode) {
	f.args = append(f.args, arg)
}

type ArrIndNode struct {
	indexer     TreeNode
	arrProvider TreeNode
}

func (a *ArrIndNode) String(ind string) string {
	ret := ind + " index\n"
	ret += a.indexer.String(ind + TAB)
	ret += "\n"
	ret += a.arrProvider.String(ind + TAB)
	return ret
}

func (f *ArrIndNode) Type() TreeNodeType {
	return ARR_IND
}

func NewArrIndNode(indexer TreeNode, arrProvider TreeNode) *ArrIndNode {
	return &ArrIndNode{indexer, arrProvider}
}

// block related nodes

type StatementsContainer interface {
	AddChild(child TreeNode)
	String(ind string) string
}

type ScopeNode struct {
	children []TreeNode
}

func MakeScopeNode() *ScopeNode {
	return &ScopeNode{make([]TreeNode, 0)}
}

func (s *ScopeNode) String(ind string) string {
	ret := ind + "scope\n" 
	for _, child := range s.children {
		ret += child.String(ind+TAB) + "\n"
	}
	return globals.RandomColor(ret)
}

func (s *ScopeNode) Type() TreeNodeType {
	return SCOPE
}

func (s *ScopeNode) AddChild(child TreeNode) {
	s.children = append(s.children, child)
}

type SourceFileNode struct {
	fileName string
	filePath string
	children []TreeNode
	//todo: store exports of this file
}

func MakeSourceFileNode() *SourceFileNode {
	return &SourceFileNode{children: make([]TreeNode, 0)}
}

func (s *SourceFileNode) String(ind string) string {
	ret := ind + "source file\n"
	for _, child := range s.children {
		ret += child.String(ind+TAB) + "\n"
	}
	return ret
}

func (s *SourceFileNode) Type() TreeNodeType {
	return SCOPE
}

func (s *SourceFileNode) AddChild(child TreeNode) {
	s.children = append(s.children, child)
}

type FuncArg struct {
	key   string
	value string
}
type FuncNode struct {
	name    string
	argList []FuncArg
	scope   *ScopeNode
}

func (f *FuncNode) String(ind string) string {
	ret := ind + "function\n" + ind + "\n" + ind + "  args:\n"
	for i := range f.argList {
		ret += ind + TAB + f.argList[i].key + " " + f.argList[i].value + "\n"
	}
	ret += f.scope.String(ind + TAB)
	return ret
}

func (f *FuncNode) Type() TreeNodeType {
	return FUNCTION
}

func MakeFunctionNode(name string) *FuncNode {
	return &FuncNode{name, make([]FuncArg, 0), nil}
}

func (f *FuncNode) AddArg(key string, value string) {
	f.argList = append(f.argList, FuncArg{key, value})
}

type VariableDeclarationNode struct {
	declarations []TreeNode
}

func MakeVariableDeclarationNode() *VariableDeclarationNode {
	return &VariableDeclarationNode{make([]TreeNode, 0)}
}

func (v *VariableDeclarationNode) String(ind string) string {
	ret := ind + "var decl\n"
	for _, decl := range v.declarations {
		ret += decl.String(ind+TAB) + "\n"
	}
	return ret
}

func (v *VariableDeclarationNode) Type() TreeNodeType {
	return VAR_DECL
}

func (v *VariableDeclarationNode) AddDeclaration(decl TreeNode) {
	v.declarations = append(v.declarations, decl)
}

type ReturnNode struct {
	value TreeNode
}

func (r *ReturnNode) String(ind string) string {
	return ind + "return\n" + r.value.String(ind+TAB)
}

func (r *ReturnNode) Type() TreeNodeType {
	return RETURN
}

func MakeReturnNode(value TreeNode) *ReturnNode {
	return &ReturnNode{value}
}

type IfNode struct {
	condition TreeNode
	ifScope   TreeNode
	elseScope TreeNode
}

func (i *IfNode) String(ind string) string {
	return ind + "if\n" + i.condition.String(ind+TAB) + "\n" + ind + "then\n" + i.ifScope.String(ind+TAB) + ind + "else\n" + i.elseScope.String(ind+TAB)
}

func (i *IfNode) Type() TreeNodeType {
	return CONDITIONAL
}

func MakeIfNode(condition TreeNode, ifScope TreeNode, elseScope TreeNode) *IfNode {
	return &IfNode{condition, ifScope, elseScope}
}

type LoopNode struct {
	initializer TreeNode
	condition   TreeNode
	updater     TreeNode
	scope       *ScopeNode
}

func (l *LoopNode) String(ind string) string {
	ret := ind + "loop\n"
	ret += l.initializer.String(ind + TAB) + "\n"
	ret += l.condition.String(ind + TAB) + "\n"
	ret += l.updater.String(ind + TAB) + "\n"
	ret += l.scope.String(ind + TAB) + "\n"
	return ret
}

func (l *LoopNode) Type() TreeNodeType {
	return LOOP
}

func MakeLoopNode(initializer TreeNode, condition TreeNode, updater TreeNode, scope *ScopeNode) *LoopNode {
	return &LoopNode{initializer, condition, updater, scope}
}
