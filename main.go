package main

import (
	"fmt"
	"he++/asm_gen"
	cmdlineutils "he++/cmdline_utils"
	"he++/lexer"
	"he++/parser"
	staticanalyzer "he++/static_analyzer"
	"he++/tac"
	"he++/utils"
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
	if !ok {
		println("Cannot proceed due to these errors")
		return
	}
	tac := tac.NewTACGen(node)
	tac.GenerateTac()
	asm_gen := asm_gen.NewAsmGen(tac)
	asm_gen.GenerateAsm()
}
