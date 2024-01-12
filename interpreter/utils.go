package interpreter

import (
	"fmt"
	"he++/parser"
	"he++/utils"
	"math"
	"os"
	"strings"

	"github.com/gofrs/uuid"
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
	fmt.Print(utils.BGCols["YELLOW"])
	fmt.Println(k...)
	fmt.Print(utils.Colors["RESET"])
}

func printVariableList(variables map[string]Variable) {
	for k, v := range variables {
		debug_info(k, v.pointer, getNumber(v))
	}
}

// todo: move to parsing phase for faster execution
func StringToNumber(str string) float64 {
	base := 10
	num := ""
	if len(str) < 2 {
		num = str
	} else {

		switch str[0:2] {
		case "0x":
			base = 16
			num = str[2:]
		case "0b":
			base = 2
			num = str[2:]
		case "0o":
			base = 8
			num = str[2:]
		default:
			num = str
		}
	}
	parsedNum := 0.0
	dotsep := strings.Split(num, ".")
	if len(dotsep) > 2 {
		interrupt("invalid number " + str)
	} else if len(dotsep) == 2 {
		parsedNum = parseNumber(dotsep[0], base) + parseFraction(dotsep[1], base)
	} else if len(dotsep) == 1 {
		parsedNum = parseNumber(dotsep[0], base)
	} else {
		interrupt("invalid number " + str)
	}
	return parsedNum
}

func parseNumber(num string, base int) float64 {
	parsedNum := 0.0
	l := len(num)
	for i := 0; i < l; i++ {
		parsedNum += float64(numVal(num[i])) * math.Pow(float64(base), float64(l-i-1))
	}
	return parsedNum
}
func parseFraction(num string, base int) float64 {
	parsedNum := 0.0
	l := len(num)
	for i := 0; i < l; i++ {
		parsedNum += float64(numVal(num[i])) * math.Pow(float64(base), float64(-(i+1)))
	}
	return parsedNum
}
func numVal(c byte) int {
	if c >= '0' && c <= '9' {
		return int(c - '0')
	}
	if c >= 'a' && c <= 'z' {
		return int(c - 'a' + 10)
	}
	if c >= 'A' && c <= 'Z' {
		return int(c - 'A' + 10)
	}
	interrupt("invalid number")
	return -1
}

func generateId() string {
	return uuid.Must(uuid.NewV4()).String()
}

func isCompositeDS(node parser.TreeNode) bool {
	return len(node.Children) > 0 || len(node.Properties) > 0
}
