package node_types

import (
	"fmt"
	"he++/utils"
)

type FuncArg struct {
	Name  string
	DataT DataType
}
type FuncNode struct {
	Name       string
	ArgList    []FuncArg
	Scope      *ScopeNode
	ReturnType DataType
	NodeMetadata
}

func (f *FuncNode) String(ind string) string {
	ret := fmt.Sprintf("%sfunc %s \n %s args:\n", ind, utils.Green(f.Name), ind)
	for i := range f.ArgList {
		ret += ind + TAB + utils.Green(f.ArgList[i].Name) + " " + utils.Cyan(f.ArgList[i].DataT.Text()) + "\n"
	}
	ret += f.Scope.String(ind + TAB)
	ret += ind + "  return type: " + utils.Cyan(f.ReturnType.Text())
	return ret
}

func (f *FuncNode) Type() TreeNodeType {
	return FUNCTION
}

func MakeFunctionNode(name string, args []FuncArg, dt DataType, scp *ScopeNode, meta *NodeMetadata) *FuncNode {
	return &FuncNode{Name: name, ArgList: args, Scope: scp, ReturnType: dt, NodeMetadata: *meta}
}

type ReturnNode struct {
	Value TreeNode
	NodeMetadata
}

func (r *ReturnNode) String(ind string) string {
	return ind + "return\n" + r.Value.String(ind+TAB)
}

func (r *ReturnNode) Type() TreeNodeType {
	return RETURN
}

func MakeReturnNode(value TreeNode, meta *NodeMetadata) *ReturnNode {
	return &ReturnNode{value, *meta}
}

type FuncCallNode struct {
	Callee TreeNode
	Args   []TreeNode
	NodeMetadata
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

func NewFuncCallNode(name TreeNode, meta *NodeMetadata) *FuncCallNode {
	return &FuncCallNode{name, make([]TreeNode, 0), *meta}
}
