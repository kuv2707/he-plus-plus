package main

import (
	"fmt"
	"toylingo/interpreter"
	"toylingo/lexer"
	"toylingo/parser"
)

func main() {

	// 	fmt.Println(utils.ValidVariableName("`aaa`"))
	// return
	var tokens *lexer.Node = lexer.Lexify("./samples/sample.lg")
	tokens = tokens.Next
	// PrintLexemes(tokens)

	treeNode := parser.ParseTreeM(tokens)

	treeNode.PrintTree("")
	fmt.Println("starting execution")
	interpreter.Interpret(treeNode)
	fmt.Println("program executed")

}
func PrintLexemes(tokens *lexer.Node) {
	c := 0
	for node := tokens; node != nil; node = node.Next {
		fmt.Println(c, node.Val.Type, node.Val.Ref)
		c++
	}
}
