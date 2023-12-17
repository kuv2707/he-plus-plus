package interpreter

import (
	"fmt"
	"os"
	"toylingo/utils"
)

func debug_error(k ...interface{}) {
	if os.Getenv("DEBUG_ERROR") == "false" {
		return
	}
	fmt.Print(utils.Colors["RED"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}

func debug_info(k ...interface{}) {
	if os.Getenv("DEBUG_INFO") == "false" {
		return
	}
	fmt.Print(utils.Colors["BOLDGREEN"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}
