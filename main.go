package main

import (
	"fmt"
	"toylingo/interpreter"
	"toylingo/lexer"
	"toylingo/parser"
	"toylingo/utils"
)

func main() {

	// 	fmt.Println(utils.ValidVariableName("`aaa`"))
	// return
	var tokens *lexer.Node = lexer.Lexify("./samples/sample.lg")
	tokens = tokens.Next
	// PrintLexemes(tokens)

	treeNode := parser.ParseTree(tokens)

	treeNode.PrintTree("")
	StartInterpreting(treeNode)
	fmt.Println(utils.Colors["RESET"])
}
func PrintLexemes(tokens *lexer.Node) {
	c := 0
	for node := tokens; node != nil; node = node.Next {
		fmt.Println(c, node.Val.Type, node.Val.Ref)
		c++
	}
}

func StartInterpreting(treeNode *parser.TreeNode) {
	fmt.Println(utils.Colors["BOLDYELLOW"] + "starting execution" + utils.Colors["RESET"])
	interpreter.Interpret(treeNode)
	fmt.Println("\n" + utils.Colors["BOLDYELLOW"] + "execution complete" + utils.Colors["RESET"])
}
