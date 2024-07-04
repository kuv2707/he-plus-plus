package utils

import (
	"fmt"
)

var indentation string = ""

const ONETAB = "    "
const endl = "\n"

func PushIndent() {
	indentation += ONETAB
}

func PopIndent() {
	if len(indentation) >= len(ONETAB) {
		indentation = indentation[:len(indentation)-len(ONETAB)]
	}
}

func IndentLog(s string) {
	Log(indentation + s)
}

func Log(s string) {
	fmt.Print(s)
}

func Logln(s string) {
	fmt.Println(s)
}

func Red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func Green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func Yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}

func Blue(s string) string {
	return "\033[34m" + s + "\033[0m"
}

func Magenta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func Cyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}

func White(s string) string {
	return "\033[37m" + s + "\033[0m"
}

func Bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func Underline(s string) string {
	return "\033[4m" + s + "\033[0m"
}

func Reverse(s string) string {
	return "\033[7m" + s + "\033[0m"
}

func Clear() {
	fmt.Print("\033[2J")
}
