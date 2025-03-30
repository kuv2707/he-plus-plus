package node_types

import "fmt"

type FuncArg struct {
	Name      string
	DataT DataType
}
type FuncNode struct {
	Name       string
	ArgList    []FuncArg
	Scope      *ScopeNode
	ReturnType DataType
}

func (f *FuncNode) String(ind string) string {
	ret := fmt.Sprintf("%sfunc %s \n %s args:\n", ind, f.Name, ind)
	for i := range f.ArgList {
		ret += ind + TAB + f.ArgList[i].Name + " " + f.ArgList[i].DataT.Text() + "\n"
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
	f.ArgList = append(f.ArgList, FuncArg{name, datatype})
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
	Callee TreeNode
	Args   []TreeNode
}

func (f *FuncCallNode) String(ind string) string {
	args := ""
	for _, arg := range f.Args {
		args += arg.String(ind+TAB+"  ") + "\n"
	}
	return ind + "call\n" + f.Callee.String(ind+TAB) + ind + "\n" + ind + "  args:\n" + args
}

func (f *FuncCallNode) Type() TreeNodeType {
	return FUNCTION_CALL
}

func NewFuncCallNode(name TreeNode) *FuncCallNode {
	return &FuncCallNode{name, make([]TreeNode, 0)}
}

func (f *FuncCallNode) Arg(arg TreeNode) {
	f.Args = append(f.Args, arg)
}
