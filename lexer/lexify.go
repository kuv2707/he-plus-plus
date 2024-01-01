package lexer

import (
	"bytes"
	"os"
	"strings"
	g "toylingo/globals"
	"toylingo/utils"
)

func Lexify(path string) *Node {

	filecontent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	//replace comments with spaces
	for i := 0; i < len(filecontent); i++ {
		if string(filecontent[i:i+2]) == "//" {
			for j := i; j < len(filecontent); j++ {
				if filecontent[j] == '\n' {
					i = j
					break
				}
				filecontent[j] = ' '
			}
		}
	}

	//decompose template literals

	// fmt.Println(string(filecontent))
	stringliterals := make([]string, 0)
	//placeholder for strings
	for i := 0; i < len(filecontent); i++ {
		if filecontent[i] == '`' {
			for j := i + 1; j < len(filecontent); j++ {
				if filecontent[j] == '`' {
					str := string(filecontent[i : j+1])

					stringliterals = append(stringliterals, str)
					i = j + 1
					break
				}
			}
		}
	}
	//parse escape sequence from all string literals
	// for i := 0; i < len(stringliterals); i++ {
	// 	stringliterals[i] = utils.ParseEscapeSequence(stringliterals[i])

	// }
	for i := 0; i < len(stringliterals); i++ {
		filecontent = bytes.ReplaceAll(filecontent, []byte(stringliterals[i]), []byte(" __STR__ "))
	}

	toPad := [...]string{"{", "}", ";", ":", "(", ")", ".", "=", "*", "/", "+", "-", "<", ">", "!", "|", "&", ",", "[", "]"}
	for i := 0; i < len(toPad); i++ {
		filecontent = bytes.ReplaceAll(filecontent, []byte(toPad[i]), []byte(" "+toPad[i]+" "))
	}
	filecontent = append(filecontent, []byte(" ")...)

	// fmt.Println(string(filecontent))
	// fmt.Println(stringliterals)
	tokens := &Node{TokenType{"start", "",0}, nil}
	ret := tokens
	temp := ""
	for i := 0; i < len(filecontent); i++ {
		c := filecontent[i]
		if c == ' ' || c == '\n' || c == '\t' {
			if c == '\n' {
				lineNo++
			}
			if addToken(strings.Trim(temp, " "), tokens) {
				tokens = tokens.Next
			}
			temp = ""
			continue
		}
		temp += string(c)

	}

	//coalesce decimal nos (4.67 etc) into one
	for node := ret; node != nil; node = node.Next {
		if node.Val.Type == "NUMBER" && node.Next != nil && node.Next.Val.Type == "DOT" && node.Next.Next != nil && node.Next.Next.Val.Type == "NUMBER" {
			node.Val.Ref = node.Val.Ref + "." + node.Next.Next.Val.Ref
			node.Next = node.Next.Next.Next
		}
	}

	//coalesce multicharacter operators into one
	for node := ret; node != nil; node = node.Next {
		if (utils.IsOneOf(node.Val.Ref, []string{"<", ">", "=", "!"}) && node.Next != nil && node.Next.Val.Ref == "=") || (node.Val.Ref == "|" && node.Next != nil && node.Next.Val.Ref == "|") || (node.Val.Ref == "&" && node.Next != nil && node.Next.Val.Ref == "&") {
			node.Val.Ref = node.Val.Ref + node.Next.Val.Ref
			node.Next = node.Next.Next
		}
		if node.Val.Ref == "+" && node.Next != nil && node.Next.Val.Ref == "+" {
			node.Val.Ref = "++"
			node.Next = node.Next.Next
		}
		if node.Val.Ref == "-" && node.Next != nil && node.Next.Val.Ref == "-" {
			node.Val.Ref = "--"
			node.Next = node.Next.Next
		}
	}

	//replace placeholder for strings
	count := 0
	for node := ret; node != nil; node = node.Next {
		if node.Val.Type == "IDENTIFIER" && node.Val.Ref == "__STR__" {
			node.Val.Ref = utils.ParseEscapeSequence(stringliterals[count])
			count++
			node.Val.Type = "STRING_LITERAL"
		}
	}

	return ret
}

var lineNo = 1

func addToken(temp string, tokens *Node) bool {
	if temp == " " || temp == "\n" || temp == "\t" || temp == "" || temp == "\r" {
		return false
	}
	switch strings.Trim(temp, " ") {
	case g.SCOPE_START:
		tokens.Next = &Node{TokenType{"SCOPE_START", g.SCOPE_START, lineNo}, nil}
	case g.SCOPE_END:
		tokens.Next = &Node{TokenType{"SCOPE_END", g.SCOPE_END, lineNo}, nil}
	case g.OPEN_PAREN:
		tokens.Next = &Node{TokenType{"OPEN_PAREN", g.OPEN_PAREN, lineNo}, nil}
	case g.CLOSE_PAREN:
		tokens.Next = &Node{TokenType{"CLOSE_PAREN", g.CLOSE_PAREN, lineNo}, nil}
	case g.OPEN_SQUARE:
		tokens.Next = &Node{TokenType{"OPEN_SQUARE", g.OPEN_SQUARE, lineNo}, nil}
	case g.CLOSE_SQUARE:
		tokens.Next = &Node{TokenType{"CLOSE_SQUARE", g.CLOSE_SQUARE, lineNo}, nil}
	case g.COLON:
		tokens.Next = &Node{TokenType{"COLON", g.COLON, lineNo}, nil}
	case g.SEMICOLON:
		tokens.Next = &Node{TokenType{"SEMICOLON", g.SEMICOLON, lineNo}, nil}
	case g.DOT:
		tokens.Next = &Node{TokenType{"DOT", g.DOT, lineNo}, nil}
	case g.LET:
		tokens.Next = &Node{TokenType{"LET", g.LET, lineNo}, nil}
	case g.IF:
		tokens.Next = &Node{TokenType{"IF", g.IF, lineNo}, nil}
	case g.ELSE_IF:
		tokens.Next = &Node{TokenType{"ELSE IF", g.ELSE_IF, lineNo}, nil}
	case g.ELSE:
		tokens.Next = &Node{TokenType{"ELSE", g.ELSE, lineNo}, nil}
	case g.LOOP:
		tokens.Next = &Node{TokenType{"LOOP", g.LOOP, lineNo}, nil}
	case g.BREAK:
		tokens.Next = &Node{TokenType{"BREAK", g.BREAK, lineNo}, nil}
	case g.COMMA:
		tokens.Next = &Node{TokenType{"COMMA", g.COMMA, lineNo}, nil}
	case g.RETURN:
		tokens.Next = &Node{TokenType{"RETURN", g.RETURN, lineNo}, nil}
	case g.FUNCTION:
		tokens.Next = &Node{TokenType{"FUNCTION", g.FUNCTION, lineNo}, nil}

	default:
		if utils.IsNumber(temp) {
			tokens.Next = &Node{TokenType{"NUMBER", temp, lineNo}, nil}
		} else if utils.IsBoolean(temp) {
			tokens.Next = &Node{TokenType{"BOOLEAN", temp, lineNo}, nil}
		} else if utils.IsOperator(temp) {
			tokens.Next = &Node{TokenType{"OPERATOR", temp, lineNo}, nil}
		} else {
			tokens.Next = &Node{TokenType{"IDENTIFIER", temp, lineNo}, nil}
		}

	}
	return true
}
