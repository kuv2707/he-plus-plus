package main

import (
	"fmt"
	cmdlineutils "he++/cmdline_utils"
	"he++/lexer"
	"he++/parser"
	staticanalyzer "he++/static_analyzer"
	"he++/utils"
	"os"
)

// "runtime/pprof"

func main() {
	args := cmdlineutils.ReadArgs()
	lexer := lexer.LexerOf(string(utils.ReadFileContent(args["src"])))
	lexer.Lexify()
	
	if os.Getenv("DEBUG_LEXER") == "1" {
		lexer.PrintLexemes()
	}
	
	astParser := parser.NewParser(lexer.GetTokens())
	node := astParser.ParseAST()
	
	if os.Getenv("DEBUG_AST") == "1" {
		fmt.Println(node.String(""))
	}
	analyzer := staticanalyzer.MakeAnalyzer()
	for _,k := range analyzer.AnalyzeAST(node) {
		fmt.Println(k)
	}
}
