package interpreter

import (
	"fmt"
	"os"
	"toylingo/utils"
)

func debug_error(k ...interface{}) {
	if os.Getenv("DEBUG_ERROR") == "0" {
		return
	}
	fmt.Print(utils.Colors["RED"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}

func debug_info(k ...interface{}) {
	if os.Getenv("DEBUG_INFO") == "0" {
		return
	}
	fmt.Print(utils.Colors["BOLDGREEN"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}

func debug_warn(k ...interface{}) {
	if os.Getenv("DEBUG_WARN") == "0" {
		return
	}
	fmt.Print(utils.Colors["YELLOW"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}

func interrupt(k ...interface{}) {
	fmt.Print(utils.Colors["RED"])
	fmt.Println(k...)
	
	fmt.Print(utils.Colors["RESET"])
	fmt.Print(utils.Colors["BOLDRED"])
	fmt.Println("execution interrupted")
	fmt.Print(utils.Colors["RESET"])
	os.Exit(1)
}