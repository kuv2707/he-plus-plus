package utils

import (
	"fmt"
	"strings"
)

var ONETAB = "  "

type ASTPrinter struct {
	OneTab  string
	indents int
	Builder strings.Builder
}

func MakeASTPrinter() ASTPrinter {
	return ASTPrinter{OneTab: ONETAB, indents: -1, Builder: strings.Builder{}}
}

func (p *ASTPrinter) PushIndent() {
	p.indents += 1
}

func (p *ASTPrinter) PopIndent() {
	p.indents -= 1
}

func (p *ASTPrinter) WriteLine(s string) {
	for range p.indents {
		p.Builder.WriteString(p.OneTab)
	}
	p.Builder.WriteString(s)
	p.Builder.WriteByte('\n')
}

func Log(s string) {
	fmt.Print(s)
}

func Logln(s string) {
	fmt.Println(s)
}

// --- Core wrapper ---
func wrap(code, s string) string {
	return code + s + "\033[0m"
}

// escape character is 0x1b or 033

// --- Foreground colors ---
func Black(s string) string   { return wrap("\033[30m", s) }
func Red(s string) string     { return wrap("\033[31m", s) }
func Green(s string) string   { return wrap("\033[32m", s) }
func Yellow(s string) string  { return wrap("\033[33m", s) }
func Blue(s string) string    { return wrap("\033[34m", s) }
func Magenta(s string) string { return wrap("\033[35m", s) }
func Cyan(s string) string    { return wrap("\033[36m", s) }
func White(s string) string   { return wrap("\033[37m", s) }

func BrightBlack(s string) string   { return wrap("\033[90m", s) }
func BrightRed(s string) string     { return wrap("\033[91m", s) }
func BrightGreen(s string) string   { return wrap("\033[92m", s) }
func BrightYellow(s string) string  { return wrap("\033[93m", s) }
func BrightBlue(s string) string    { return wrap("\033[94m", s) }
func BrightMagenta(s string) string { return wrap("\033[95m", s) }
func BrightCyan(s string) string    { return wrap("\033[96m", s) }
func BrightWhite(s string) string   { return wrap("\033[97m", s) }

func BoldBlack(s string) string   { return wrap("\033[1m\033[30m", s) }
func BoldRed(s string) string     { return wrap("\033[1m\033[31m", s) }
func BoldGreen(s string) string   { return wrap("\033[1m\033[32m", s) }
func BoldYellow(s string) string  { return wrap("\033[1m\033[33m", s) }
func BoldBlue(s string) string    { return wrap("\033[1m\033[34m", s) }
func BoldMagenta(s string) string { return wrap("\033[1m\033[35m", s) }
func BoldCyan(s string) string    { return wrap("\033[1m\033[36m", s) }
func BoldWhite(s string) string   { return wrap("\033[1m\033[37m", s) }

func BgBlack(s string) string   { return wrap("\033[40m", s) }
func BgRed(s string) string     { return wrap("\033[41m", s) }
func BgGreen(s string) string   { return wrap("\033[42m", s) }
func BgYellow(s string) string  { return wrap("\033[43m", s) }
func BgBlue(s string) string    { return wrap("\033[44m", s) }
func BgMagenta(s string) string { return wrap("\033[45m", s) }
func BgCyan(s string) string    { return wrap("\033[46m", s) }
func BgWhite(s string) string   { return wrap("\033[47m", s) }

func BgBoldBlack(s string) string   { return wrap("\033[1m\033[40m", s) }
func BgBoldRed(s string) string     { return wrap("\033[1m\033[41m", s) }
func BgBoldGreen(s string) string   { return wrap("\033[1m\033[42m", s) }
func BgBoldYellow(s string) string  { return wrap("\033[1m\033[43m", s) }
func BgBoldBlue(s string) string    { return wrap("\033[1m\033[44m", s) }
func BgBoldMagenta(s string) string { return wrap("\033[1m\033[45m", s) }
func BgBoldCyan(s string) string    { return wrap("\033[1m\033[46m", s) }
func BgBoldWhite(s string) string   { return wrap("\033[1m\033[47m", s) }

// --- Text styles ---
func Bold(s string) string      { return wrap("\033[1m", s) }
func Underline(s string) string { return wrap("\033[4m", s) }
func Reverse(s string) string   { return wrap("\033[7m", s) }
func Reset(s string) string     { return wrap("\033[0m", s) }
