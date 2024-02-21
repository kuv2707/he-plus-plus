package globals

var SCOPE_START = "{"
var SCOPE_END = "}"
var COLON = ":"
var SEMICOLON = ";"
var LET = "let"
var INTEGER = "int"
var FLOAT = "float"
var BOOLEAN = "bool"
var STRING = "string"
var DOT = "."
var OPEN_PAREN = "("
var CLOSE_PAREN = ")"
var IF = "if"
var ELSE_IF = "elseif"
var ELSE = "else"
var LOOP = "loop"
var BREAK = "break"
var FUNCTION = "function"
var COMMA = ","
var RETURN = "return"
var EQUALS = "="
var OPEN_SQUARE = "["
var CLOSE_SQUARE = "]"

var KEYWORDS = []string{"IF", "ELSE IF", "ELSE", "FUNCTION", "SCOPE_END", "LET"}

var Operators = []string{"=", "+", "-", "*", "/", "<", ">", "#", "!", "|", "&", "==", "!=", "<=", ">=", "++", "--", "&&", "||"}

//todo bitwise operators

var Booleans = []string{"true", "false"}

//contains mapping of number strings to number values and strings to their hash values
var NumMap = map[string][]byte{}


func HashString(s string) int64 {
	var val int64 = 0
	for i := 0; i < len(s); i++ {
		val += int64(s[i]) * 51
	}
	return val
}
