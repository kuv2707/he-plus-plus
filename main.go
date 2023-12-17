package main

import (
	"fmt"
	"os"
	"toylingo/interpreter"
	"toylingo/lexer"
	"toylingo/parser"
	"toylingo/utils"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file")
	}
	var tokens *lexer.Node = lexer.Lexify("./samples/"+os.Getenv("SOURCE_FILE"))
	tokens = tokens.Next
	if os.Getenv("DEBUG_LEXER") == "true" {
		PrintLexemes(tokens)
	}
	treeNode := parser.ParseTree(tokens)

	if os.Getenv("DEBUG_AST") == "true" {
		treeNode.PrintTree("")
	}
	if os.Getenv("INTERPRET") == "true" {
		StartInterpreting(treeNode)
	}
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
