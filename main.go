package main

import(
	"fmt"
	"toylingo/lexer"
	// "toylingo/parser"
)

func main(){

	var tokens *lexer.Node =lexer.Lexify("./samples/sample.lg")
	for node:=tokens;node!=nil;node=node.Next{
		fmt.Print(node.Val)
	}
	
	// treeNode:=parser.ParseTree(tokens)

	// parser.PrintTree(treeNode,"")



}