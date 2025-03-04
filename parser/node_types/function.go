package node_types

import "fmt"

type FuncArg struct {
	name      string
	data_type DataType
}
type FuncNode struct {
	name       string
	argList    []FuncArg
	Scope      *ScopeNode
	ReturnType string
}

func (f *FuncNode) String(ind string) string {
	ret := fmt.Sprintf("%sfunc %s \n %s args:\n", ind, f.name, ind)
	for i := range f.argList {
		ret += ind + TAB + f.argList[i].name + " " + f.argList[i].data_type.Text + "\n"
	}
	ret += f.Scope.String(ind + TAB)
	ret += ind + "  return type: " + f.ReturnType
	return ret
}

func (f *FuncNode) Type() TreeNodeType {
	return FUNCTION
}

func MakeFunctionNode(name string) *FuncNode {
	return &FuncNode{name, make([]FuncArg, 0), nil, ""}
}

func (f *FuncNode) AddArg(name string, datatype DataType) {
	f.argList = append(f.argList, FuncArg{name, datatype})
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

func (f *FuncCallNode) Arg(arg TreeNode) {
	f.args = append(f.args, arg)
}
