package main

import (
	"fmt"
	cmdlineutils "he++/cmdline_utils"
	"he++/lexer"
	"he++/parser"
	staticanalyzer "he++/static_analyzer"
	"he++/utils"
	x64sysvasm "he++/x64_sysv_asm"
	"os"
)

// "runtime/pprof"

func main() {
	args := cmdlineutils.ReadArgs()
	lexer := lexer.LexerOf(args["src"])
	go lexer.Lexify()

	astParser := parser.NewParser(lexer)
	node := astParser.ParseAST()

	if os.Getenv("DEBUG_LEXER") == "1" {
		fmt.Println("Lexemes:")
		lexer.PrintLexemes()
	}
	if os.Getenv("DEBUG_AST") == "1" {
		p := utils.MakeASTPrinter()
		node.String(&p)
		fmt.Println(p.Builder.String())
	}
	analyzer := staticanalyzer.MakeAnalyzer()
	ok := analyzer.AnalyzeAST(node)
	if ok {
		asm := x64sysvasm.NewTACGen(node)
		asm.GenerateTac()
	}
}
