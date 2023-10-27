package lexer

import (
	"bytes"
	_"fmt"
	"os"
	"strings"
	"toylingo/utils"
)

func Lexify(path string) *Node {

	filecontent, err := os.ReadFile(path)
	if err != nil {
		panic(err)
	}

	toPad := [...]string{"{", "}", ";", ":", "(", ")", ".", "=", "*", "/", "+", "-", "<", ">", "!"}
	for i := 0; i < len(toPad); i++ {
		filecontent = bytes.ReplaceAll(filecontent, []byte(toPad[i]), []byte(" "+toPad[i]+" "))
	}
	filecontent = append(filecontent, []byte(" ")...)

	stringliterals := make([]string, 0)
	//placeholder for strings
	for i := 0; i < len(filecontent); i++ {
		if filecontent[i] == '`' {
			for j := i + 1; j < len(filecontent); j++ {
				if filecontent[j] == '`' {
					str := string(filecontent[i : j+1])
					stringliterals = append(stringliterals, str)
					i = j+1
					break
				}
			}
		}
	}
	for i := 0; i < len(stringliterals); i++ {
		filecontent = bytes.ReplaceAll(filecontent, []byte(stringliterals[i]), []byte(" __STR__ "))
	}

	// fmt.Println(string(filecontent))
	// fmt.Println(stringliterals)
	tokens := &Node{TokenType{"start", ""}, nil}
	ret := tokens
	temp := ""
	for i := 0; i < len(filecontent); i++ {
		c := filecontent[i]
		if c == ' ' || c == '\n' || c == '\t' {
			// fmt.Println("word ", temp)
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
		if utils.IsOneOf(node.Val.Ref, "<>=") && node.Next != nil && node.Next.Val.Ref == "=" {
			node.Val.Ref = node.Val.Ref + node.Next.Val.Ref
			node.Next = node.Next.Next
		}
	}

	//replace placeholder for strings
	count := 0
	for node := ret; node != nil; node = node.Next {
		if node.Val.Type == "IDENTIFIER" && node.Val.Ref == "__STR__" {
			node.Val.Ref = stringliterals[count]
			count++
			node.Val.Type = "STRING_LITERAL"

		}
	}

	return ret
}

func addToken(temp string, tokens *Node) bool {
	if temp == " " || temp == "\n" || temp == "\t" || temp == "" || temp == "\r" {
		return false
	}
	switch strings.Trim(temp, " ") {
	case dict["SCOPE_START"]:
		tokens.Next = &Node{TokenType{"SCOPE_START", ""}, nil}
	case dict["SCOPE_END"]:
		tokens.Next = &Node{TokenType{"SCOPE_END", ""}, nil}
	case dict["OPEN_PAREN"]:
		tokens.Next = &Node{TokenType{"OPEN_PAREN", "("}, nil}
	case dict["CLOSE_PAREN"]:
		tokens.Next = &Node{TokenType{"CLOSE_PAREN", ")"}, nil}
	case dict["COLON"]:
		tokens.Next = &Node{TokenType{"COLON", ""}, nil}
	case dict["SEMICOLON"]:
		tokens.Next = &Node{TokenType{"SEMICOLON", ""}, nil}
	case dict["LET"]:
		tokens.Next = &Node{TokenType{"LET", ""}, nil}
	case dict["IF"]:
		tokens.Next = &Node{TokenType{"IF", ""}, nil}
	case dict["DOT"]:
		tokens.Next = &Node{TokenType{"DOT", ""}, nil}

	case dict["INTEGER"]:
		tokens.Next = &Node{TokenType{"DATATYPE", "INTEGER"}, nil}
	case dict["FLOAT"]:
		tokens.Next = &Node{TokenType{"DATATYPE", "FLOAT"}, nil}
	case dict["STRING"]:
		tokens.Next = &Node{TokenType{"DATATYPE", "STRING"}, nil}
	case dict["BOOLEAN"]:
		tokens.Next = &Node{TokenType{"DATATYPE", "BOOLEAN"}, nil}
	default:
		if utils.IsNumber(temp) {
			tokens.Next = &Node{TokenType{"NUMBER", temp}, nil}
		} else if isOperator(temp) {
			tokens.Next = &Node{TokenType{"OPERATOR", temp}, nil}
		} else {
			tokens.Next = &Node{TokenType{"IDENTIFIER", temp}, nil}
		}

	}
	return true
}

func isOperator(temp string) bool {
	operators := "=+-*/<>#!"
	return utils.IsOneOf(temp, operators)
}
