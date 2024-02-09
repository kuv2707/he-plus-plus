package interpreter

import (
	"fmt"
	"he++/parser"
	"he++/utils"
	// "math"
	"os"
	// "strings"

	"github.com/gofrs/uuid"
)

var LineNo = -1

func debug_error(k ...interface{}) {
	if "0" == "0" {
		return
	}
	fmt.Print(utils.Colors["RED"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}

func debug_info(k ...interface{}) {
	if "0" == "0" {
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
	fmt.Print(utils.BGCols["YELLOW"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}


func generateId() string {
	return uuid.Must(uuid.NewV4()).String()
}

func isCompositeDS(node parser.TreeNode) bool {
	return len(node.Children) > 0 || len(node.Properties) > 0
}

func StringToNumber(s string) (float64, bool) {
	var f float64
	_, err := fmt.Sscanf(s, "%f", &f)
	if err != nil {
		return 0, false
	}
	return f, true
}