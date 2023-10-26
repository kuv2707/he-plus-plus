package main

import (
	"fmt"
	"toylingo/interpreter"
	"toylingo/lexer"
	"toylingo/parser"
)

func main(){

	var tokens *lexer.Node =lexer.Lexify("./samples/sample.lg")
	// c:=0
	// for node:=tokens;node!=nil;node=node.Next{
	// 	fmt.Println(c,node.Val.Type, node.Val.Ref)
	// 	c++
	// }
	

	
	
	treeNode:=parser.ParseTreeM(tokens)

	// treeNode.PrintTree("")
	
	interpreter.ExecuteAST(treeNode)
	fmt.Println("program executed")


}