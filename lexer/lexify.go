package lexer

import (
	"fmt"
	"he++/utils"
)

type Lexer struct {
	sourceCode string
	i          int
	lineCnt    int
	tokens     []LexerToken
	word       string
	warnings   []string
}

func LexerOf(src string) *Lexer {
	return &Lexer{sourceCode: src, i: 0, lineCnt: 1, tokens: make([]LexerToken, 0), word: "", warnings: make([]string, 0)}
}

func (l *Lexer) addWarning(err string) {
	l.warnings = append(l.warnings, err)
}

func (l *Lexer) ClearWord() {
	l.word = ""
}

func (l *Lexer) PopAndInsert(i int, token LexerToken) {
	l.tokens = append(l.tokens[:len(l.tokens)-i], token)
}

func (l *Lexer) PrintLexemes() {
	for _, token := range l.tokens {
		fmt.Println(token)
	}

	if len(l.warnings) == 0 {
		return
	}
	fmt.Println("Warnings:")
	for _, err := range l.warnings {
		fmt.Println(err)
	}
}

func (l *Lexer) makeToken(word string) LexerToken {
	if isKeyword(word) {
		return LexerToken{KEYWORD, word, l.lineCnt}
	}
	return LexerToken{IDENTIFIER, word, l.lineCnt}
}

func (l *Lexer) GetTokens() []LexerToken {
	return l.tokens
}

func (l *Lexer) lookBack(i int) LexerToken {
	if len(l.tokens) > i {
		return l.tokens[len(l.tokens)-i]
	}
	return LexerToken{}
}

func (l *Lexer) addTokenAndClearWord(token LexerToken) {
	l.tokens = append(l.tokens, token)
	l.word = ""
}

func (l *Lexer) addTokenIfCan(word string) {
	if word != "" {
		l.addTokenAndClearWord(l.makeToken(word))
	}
}

func (l *Lexer) addOperatorToken(op string) {
	l.addTokenAndClearWord(NewLexerToken(OPERATOR, op, l.lineCnt))
}

func (l *Lexer) addOperator() {
	if l.i+1 < len(l.sourceCode) && isOperator(l.sourceCode[l.i:l.i+2]) {
		l.addOperatorToken(l.sourceCode[l.i : l.i+2])
		l.i++
	} else {
		l.addOperatorToken(l.sourceCode[l.i : l.i+1])
	}
}

func (l *Lexer) escapeSequence(c byte) string {
	ret := ""
	switch c {
	case 'n':
		ret += "\n"
	case 't':
		ret += "\t"
	case 'r':
		ret += "\r"
	case 'b':
		ret += "\b"
	case 'f':
		ret += "\f"
	case '\\':
		ret += "\\"
	case '\'':
		ret += "`"
	case '"':
		ret += "\""
	default:
		l.addWarning(fmt.Sprintf("Ignored escape sequence %s at line %d", utils.Blue(fmt.Sprintf("\"\\%c\"", c)), l.lineCnt))
	}
	return ret
}

func (l *Lexer) isThisLexicalQuote() bool {
	return isQuote(string(l.sourceCode[l.i])) && (l.i == 0 || l.sourceCode[l.i-1] != '\\')
}

func (l *Lexer) Lexify() []LexerToken {

	seekTillLineEnd := func() {
		for ; l.i < len(l.sourceCode) && l.sourceCode[l.i] != '\n'; l.i++ {
		}
	}

	for ; l.i < len(l.sourceCode); l.i++ {
		c := string(l.sourceCode[l.i])
		if isDelimiter(c) {
			l.addTokenIfCan(l.word)
			if c == "\n" {
				l.lineCnt++
			}

			// strings
		} else if isPunctuation(c) {
			l.addTokenIfCan(l.word)
			l.addTokenAndClearWord(NewLexerToken(PUNCTUATION, c, l.lineCnt))
		} else if l.isThisLexicalQuote() {
			l.addTokenIfCan(l.word)
			for l.i++; l.i < len(l.sourceCode) && !l.isThisLexicalQuote(); l.i++ {
				if l.sourceCode[l.i] == '\\' {
					// escape sequence
					l.i++
					l.word += l.escapeSequence(l.sourceCode[l.i])

				} else {
					l.word += string(l.sourceCode[l.i])
				}
			}
			l.addTokenAndClearWord(NewLexerToken(STRING_LITERAL, l.word, l.lineCnt))

			// comments
		} else if c == "/" {
			l.addTokenIfCan(l.word)
			if l.i+1 < len(l.sourceCode) && l.sourceCode[l.i+1] == '/' {
				seekTillLineEnd()
			} else {
				l.addOperator()
			}

		} else if isBracket(c) {
			l.addTokenIfCan(l.word)
			l.addTokenAndClearWord(NewLexerToken(BRACKET, c, l.lineCnt))

		} else if isDigit(c) {
			l.addTokenIfCan(l.word)
			for ; l.i < len(l.sourceCode) && isDigit(string(l.sourceCode[l.i])); l.i++ {
				l.word += string(l.sourceCode[l.i])
			}
			// todo: handle hex and other characters
			if l.lookBack(1).Text() == DOT && l.lookBack(2).Type() == INTEGER {
				l.PopAndInsert(2, NewLexerToken(FLOATINGPT, l.lookBack(2).Text()+DOT+l.word, l.lineCnt))
				l.ClearWord()
			} else {
				l.addTokenAndClearWord(NewLexerToken(INTEGER, l.word, l.lineCnt))
			}
			l.i--

		} else if isOperator(c) {
			l.addTokenIfCan(l.word)
			l.addOperator()
		} else {
			l.word += c
		}
	}

	l.addTokenIfCan(l.word)

	return l.tokens
}
