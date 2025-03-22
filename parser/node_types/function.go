package node_types

import "fmt"

type FuncArg struct {
	name      string
	data_type DataType
}
type FuncNode struct {
	Name       string
	argList    []FuncArg
	Scope      *ScopeNode
	ReturnType DataType
}

func (f *FuncNode) String(ind string) string {
	ret := fmt.Sprintf("%sfunc %s \n %s args:\n", ind, f.Name, ind)
	for i := range f.argList {
		ret += ind + TAB + f.argList[i].name + " " + f.argList[i].data_type.Text() + "\n"
	}
	ret += f.Scope.String(ind + TAB)
	ret += ind + "  return type: " + f.ReturnType.Text()
	return ret
}

func (f *FuncNode) Type() TreeNodeType {
	return FUNCTION
}

func MakeFunctionNode(name string) *FuncNode {
	return &FuncNode{name, make([]FuncArg, 0), nil, &VoidType{}}
}

func (f *FuncNode) AddArg(name string, datatype DataType) {
	f.argList = append(f.argList, FuncArg{name, datatype})
}

type ReturnNode struct {
	Value TreeNode
}

func (r *ReturnNode) String(ind string) string {
	return ind + "return\n" + r.Value.String(ind+TAB)
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
	return FUNCTION_CALL
}

func NewFuncCallNode(name TreeNode) *FuncCallNode {
	return &FuncCallNode{name, make([]TreeNode, 0)}
}

func (f *FuncCallNode) Arg(arg TreeNode) {
	f.args = append(f.args, arg)
}
