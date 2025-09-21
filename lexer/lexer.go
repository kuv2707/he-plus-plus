package lexer

import (
	"fmt"
	"he++/globals"
	"strings"
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

type Warning struct {
	Msg  string
	Line int
}

type Lexer struct {
	sourceCode string
	i          int
	lineCnt    int
	// todo: use chan
	tokens   []LexerToken
	word     strings.Builder
	warnings []Warning
}

func (l *Lexer) CharAtOffset(offset int) byte {
	i := l.i + offset
	if i >= len(l.sourceCode) || i < 0 {
		return 0
	}
	return l.sourceCode[i]
}

func LexerOf(src string) *Lexer {
	return &Lexer{sourceCode: src, i: 0, lineCnt: 1, word: strings.Builder{}}
}

func (l *Lexer) addWarning(err string, line int) {
	l.warnings = append(l.warnings, Warning{Line: line, Msg: err})
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

func (l *Lexer) addTokenAndClearWord(token LexerToken) {
	l.tokens = append(l.tokens, token)
	l.word.Reset()
}

func (l *Lexer) addTokenIfCan() {
	if l.word.Len() != 0 {
		l.addTokenAndClearWord(l.makeToken(l.word.String()))
	}
}

func (l *Lexer) addOperatorToken(op string) {
	l.addTokenAndClearWord(NewLexerToken(OPERATOR, op, l.lineCnt))
}

func (l *Lexer) tryOperator() bool {
	if l.i+1 >= len(l.sourceCode) {
		return false
	}
	offset := OpTrie.MatchLongest(l.sourceCode, l.i)
	if offset != -1 {
		l.addTokenIfCan()
		l.addOperatorToken(l.sourceCode[l.i : l.i+1+offset])
		l.i += offset
	} else {
		return false
	}
	return true
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
		l.addWarning(fmt.Sprintf("Ignored escape sequence %s at line %d", globals.Blue(fmt.Sprintf("\"\\%c\"", c)), l.lineCnt), l.lineCnt)
	}
	return ret
}

func (l *Lexer) isThisLexicalQuote() bool {
	return isQuote(l.CharAtOffset(0)) && (l.CharAtOffset(-1) != '\\')
}
