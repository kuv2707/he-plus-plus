package utils
import (
	"math/rand"
)
var Colors = map[string]string{
	"RESET":         "\033[0m",
	"BLACK":         "\033[30m",
	"RED":           "\033[31m",
	"GREEN":         "\033[32m",
	"YELLOW":        "\033[33m",
	"BLUE":          "\033[34m",
	"MAGENTA":       "\033[35m",
	"CYAN":          "\033[36m",
	"WHITE":         "\033[37m",
	"BOLDBLACK":     "\033[1m\033[30m",
	"BOLDRED":       "\033[1m\033[31m",
	"BOLDGREEN":     "\033[1m\033[32m",
	"BOLDYELLOW":    "\033[1m\033[33m",
	"BOLDBLUE":      "\033[1m\033[34m",
	"BOLDMAGENTA":   "\033[1m\033[35m",
	"BOLDCYAN":      "\033[1m\033[36m",
	"BOLDWHITE":     "\033[1m\033[37m",
	
}
var BGCols = map[string]string{
	"BLACK":      "\033[40m",
	"RED":        "\033[41m",
	"GREEN":      "\033[42m",
	"YELLOW":     "\033[43m",
	"BLUE":       "\033[44m",
	"MAGENTA":    "\033[45m",
	"CYAN":       "\033[46m",
	"WHITE":      "\033[47m",
	"BOLDBLACK":  "\033[1m\033[40m",
	"BOLDRED":    "\033[1m\033[41m",
	"BOLDGREEN":  "\033[1m\033[42m",
	"BOLDYELLOW": "\033[1m\033[43m",
	"BOLDBLUE":   "\033[1m\033[44m",
	"BOLDMAGENTA": "\033[1m\033[45m",
	"BOLDCYAN":    "\033[1m\033[46m",
	"BOLDWHITE":   "\033[1m\033[47m",
}
func GetRandomColor() string {
	for _, color := range Colors {
		if color == Colors["RESET"] {
			continue
		}
		if rand.Intn(2) == 1 {

			return color
		}
	}
	return Colors["BLUE"]
}