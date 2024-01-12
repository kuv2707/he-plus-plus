package main

import (
	"fmt"
	"he++/interpreter"
	"he++/lexer"
	"he++/parser"
	"he++/utils"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file. Make sure you have a .env file in the root directory in the format specified in the .env.example file")
	}
	var tokens *lexer.Node = lexer.Lexify("./samples/" + os.Getenv("SOURCE_FILE"))
	tokens = tokens.Next
	// return
	if os.Getenv("DEBUG_LEXER") == "1" {
		PrintLexemes(tokens)
	}
	treeNode := parser.ParseTree(tokens)

	if os.Getenv("DEBUG_AST") == "1" {
		treeNode.PrintTree("")
	}
	if os.Getenv("INTERPRET") == "1" {
		fmt.Println(utils.Colors["YELLOW"] + "--> " + os.Getenv("SOURCE_FILE"))
		StartInterpreting(treeNode)
	}
	fmt.Println(utils.Colors["RESET"])
}
func PrintLexemes(tokens *lexer.Node) {
	c := 0
	for node := tokens; node != nil; node = node.Next {
		fmt.Println(c, node.Val.Type, node.Val.Ref, node.Val.LineNo)
		c++
	}
}

func StartInterpreting(treeNode *parser.TreeNode) {
	fmt.Println(utils.Colors["BOLDYELLOW"] + "starting execution" + utils.Colors["RESET"])
	startTime := time.Now().UnixMilli()
	interpreter.Interpret(treeNode)
	endTime := time.Now().UnixMilli()
	fmt.Println("\n" + utils.Colors["BOLDYELLOW"] + "execution completed in " + fmt.Sprint(endTime-startTime) + "ms" + utils.Colors["RESET"])
}
