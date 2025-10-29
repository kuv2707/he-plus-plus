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

func (f *FuncNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine(fmt.Sprintf("%s %s", utils.Underline("function"), utils.Green(f.Name)))
	p.PushIndent()
	p.WriteLine("returns: " + utils.Cyan(f.ReturnType.Text()))
	if len(f.ArgList) > 0 {
		p.WriteLine("arglist:")
		p.PushIndent()
		for i := range f.ArgList {
			p.WriteLine(fmt.Sprintf("%s %s", utils.Green(f.ArgList[i].Name), utils.Cyan(f.ArgList[i].DataT.Text())))
		}
		p.PopIndent()
	}
	p.PopIndent()
	f.Scope.String(p)
	p.PopIndent()
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

func (r *ReturnNode) String(p *utils.ASTPrinter) {
	p.PushIndent()
	p.WriteLine("return")
	r.Value.String(p)
	p.PopIndent()
}

func (r *ReturnNode) Type() TreeNodeType {
	return RETURN
}

func MakeReturnNode(value TreeNode, meta *NodeMetadata) *ReturnNode {
	return &ReturnNode{value, *meta}
}
