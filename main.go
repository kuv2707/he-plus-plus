package main

import (
	"bufio"
	"fmt"
	"he++/interpreter"
	"he++/lexer"
	"he++/parser"
	"he++/utils"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

import (
	"os"
	// "runtime/pprof"
)

func readFileContent(path string) []byte {
	filecontent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}
	return filecontent
}

func main() {
	if err := godotenv.Load(); err != nil {
		panic("Error loading .env file. Make sure you have a .env file in the root directory in the format specified in the .env.example file")
	}

	if os.Getenv("REPL") == "1" {
		fmt.Println(utils.Colors["BOLDYELLOW"] + "REPL mode" + utils.Colors["RESET"])
		startREPL()
		return
	}

	var tokens *lexer.Node = lexer.Lexify(readFileContent("./samples/" + os.Getenv("SOURCE_FILE")))
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

func StartInterpreting(treeNode *parser.TreeNode) {
	// f, err := os.Create("cpu.pprof")
	// if err != nil {
	// 	fmt.Println("could not create CPU profile: ", err)
	// 	return
	// }
	// pprof.StartCPUProfile(f)
	// defer pprof.StopCPUProfile()
	fmt.Println(utils.Colors["BOLDYELLOW"] + "starting execution" + utils.Colors["RESET"])
	startTime := time.Now().UnixMilli()
	ctx := interpreter.Init()
	interpreter.Interpret(treeNode, ctx)
	endTime := time.Now().UnixMilli()
	fmt.Println("\n" + utils.Colors["BOLDYELLOW"] + "execution completed in " + fmt.Sprint(endTime-startTime) + "ms" + utils.Colors["RESET"])
}

// read, evaluate, print, loop
func startREPL() {
	ctx := interpreter.Init()
	for {
		fmt.Print(utils.Colors["BOLDYELLOW"] + "he++> " + utils.Colors["RESET"])
		reader := bufio.NewReader(os.Stdin)
		input, _ := reader.ReadString('\n')
		if strings.TrimSpace(input) == "exit" {
			break
		}
		var tokens *lexer.Node = lexer.Lexify([]byte(input))
		tokens = tokens.Next
		treeNode := parser.ParseTree(tokens)
		if os.Getenv("INTERPRET") == "1" {
			interpreter.Interpret(treeNode, ctx)
		}
	}
	fmt.Println(utils.Colors["BOLDYELLOW"] + "exiting" + utils.Colors["RESET"])
	fmt.Println(utils.Colors["RESET"])
}

func PrintLexemes(tokens *lexer.Node) {
	c := 0
	for node := tokens; node != nil; node = node.Next {
		fmt.Println(c, node.Val.Type, node.Val.Ref, node.Val.LineNo)
		c++
	}
}
