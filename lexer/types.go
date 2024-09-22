package lexer

import (
	"fmt"
	"he++/globals"
)

type LexerTokenType string

func (m LexerTokenType) String() string {
	return string(m)
}

type LexerToken struct {
	tokenType LexerTokenType
	ref       string
	lineNo    int
}

func NewLexerToken(tokenType LexerTokenType, ref string, lineNo int) LexerToken {
	return LexerToken{tokenType, ref, lineNo}
}

func (m LexerToken) String() string {
	return fmt.Sprintf("%s %s %s", globals.Blue(string(m.tokenType)), globals.Yellow(m.ref), globals.Red(fmt.Sprint(m.lineNo)))
}

func (l LexerToken) Text() string {
	return l.ref
}

func (l LexerToken) Type() LexerTokenType {
	return l.tokenType
}

func (l LexerToken) LineNo() int {
	return l.lineNo
}

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

func (l *Lexer) clearWord() {
	l.word = ""
}

func (l *Lexer) popAndInsert(i int, token LexerToken) {
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
		l.addWarning(fmt.Sprintf("Ignored escape sequence %s at line %d", globals.Blue(fmt.Sprintf("\"\\%c\"", c)), l.lineCnt))
	}
	return ret
}

func (l *Lexer) isThisLexicalQuote() bool {
	return isQuote(string(l.sourceCode[l.i])) && (l.i == 0 || l.sourceCode[l.i-1] != '\\')
}
