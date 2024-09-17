package main

import (
	"fmt"
	"os"
	cmdlineutils "he++/cmdline_utils"
	"he++/parser"
	"he++/lexer"
	"he++/utils"
)

// "runtime/pprof"

func main() {
	args := cmdlineutils.ReadArgs()
	lexer := lexer.LexerOf(string(utils.ReadFileContent(args["src"])))
	lexer.Lexify()
	
	astParser := parser.NewParser(lexer.GetTokens())
	node := astParser.ParseAST()
	
	if os.Getenv("DEBUG_LEXER") == "1" {
		lexer.PrintLexemes()
	}
	if os.Getenv("DEBUG_AST") == "1" {
		fmt.Println(node.String(""))
	}
	if os.Getenv("INTERPRET") == "1" {
		fmt.Println(utils.Colors["YELLOW"] + "--> " + os.Getenv("SOURCE_FILE"))
		// StartInterpreting(treeNode, args)
	}
	fmt.Println(utils.Colors["RESET"])
}

// func StartInterpreting(treeNode *parser.TreeNode, args map[string]string) {
// 	// f, err := os.Create("cpu.pprof")
// 	// if err != nil {
// 	// 	fmt.Println("could not create CPU profile: ", err)
// 	// 	return
// 	// }
// 	// pprof.StartCPUProfile(f)
// 	// defer pprof.StopCPUProfile()
// 	fmt.Println(utils.Colors["BOLDYELLOW"] + "starting execution" + utils.Colors["RESET"])
// 	startTime := time.Now().UnixMilli()
// 	ctx := interpreter.Init(args)
// 	interpreter.Interpret(treeNode, ctx)
// 	endTime := time.Now().UnixMilli()
// 	fmt.Println("\n" + utils.Colors["BOLDYELLOW"] + "execution completed in " + fmt.Sprint(endTime-startTime) + "ms" + utils.Colors["RESET"])
// }

// read, evaluate, print, loop
// func startREPL() {
// 	ctx := interpreter.Init()
// 	for {
// 		fmt.Print(utils.Colors["BOLDYELLOW"] + "he++> " + utils.Colors["RESET"])
// 		reader := bufio.NewReader(os.Stdin)
// 		input, _ := reader.ReadString('\n')
// 		if strings.TrimSpace(input) == "exit" {
// 			break
// 		}
// 		var tokens *lexer.Node = lexer.Lexify([]byte(input))
// 		tokens = tokens.Next
// 		treeNode := parser.ParseTree(tokens)
// 		if os.Getenv("INTERPRET") == "1" {
// 			interpreter.Interpret(treeNode, ctx)
// 		}
// 	}
// 	fmt.Println(utils.Colors["BOLDYELLOW"] + "exiting" + utils.Colors["RESET"])
// 	fmt.Println(utils.Colors["RESET"])
// }
