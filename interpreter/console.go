package interpreter

import (
	"fmt"
)

var indentation string = ""

const ONETAB = "    "
const endl = "\n"

func pushIndent() {
	indentation += ONETAB
}

func popIndent() {
    if len(indentation) >= len(ONETAB) {
        indentation = indentation[:len(indentation)-len(ONETAB)]
    }
}
func indentLog(s string) {
	log(indentation + s)
}

func log(s string) {
	fmt.Print(s)
}

func logln(s string) {
	fmt.Println(s)
}

func red(s string) string {
	return "\033[31m" + s + "\033[0m"
}

func green(s string) string {
	return "\033[32m" + s + "\033[0m"
}

func yellow(s string) string {
	return "\033[33m" + s + "\033[0m"
}

func blue(s string) string {
	return "\033[34m" + s + "\033[0m"
}

func magenta(s string) string {
	return "\033[35m" + s + "\033[0m"
}

func cyan(s string) string {
	return "\033[36m" + s + "\033[0m"
}

func white(s string) string {
	return "\033[37m" + s + "\033[0m"
}

func bold(s string) string {
	return "\033[1m" + s + "\033[0m"
}

func underline(s string) string {
	return "\033[4m" + s + "\033[0m"
}

func reverse(s string) string {
	return "\033[7m" + s + "\033[0m"
}

func clear() {
	fmt.Print("\033[2J")
}
