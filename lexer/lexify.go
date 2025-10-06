package lexer

import (
	"strings"
)

func (l *Lexer) Lexify() {
	for ; l.i < len(l.sourceCode); l.i++ {
		c := l.CharAtOffset(0)
		sc := string(c)
		if isDelimiter(c) {
			l.addTokenIfCan()
			if c == '\n' {
				l.lineCnt++
			}
		} else if isPunctuation(sc) {
			l.addTokenIfCan()
			l.addTokenAndClearWord(NewLexerToken(PUNCTUATION, sc, l.lineCnt))
		} else if l.isThisLexicalQuote() {
			l.addTokenIfCan()
			// strings
			for l.i++; l.i < len(l.sourceCode) && !l.isThisLexicalQuote(); l.i++ {
				if l.CharAtOffset(0) == '\\' {
					// escape sequence
					l.i++
					l.word.WriteString(l.escapeSequence(l.sourceCode[l.i]))
				} else {
					l.word.WriteByte(l.sourceCode[l.i])
				}
			}
			l.addTokenAndClearWord(NewLexerToken(STRING_LITERAL, l.word.String(), l.lineCnt))

		} else if c == '/' {
			l.addTokenIfCan()
			// comments
			if l.CharAtOffset(1) == '/' {
				for ; l.CharAtOffset(0) != '\n'; l.i++ {
				}
				l.lineCnt++
			} else {
				l.tryOperator()
			}

		} else if isBracket(sc) {
			l.addTokenIfCan()
			l.addTokenAndClearWord(NewLexerToken(BRACKET, sc, l.lineCnt))

		} else if isDigit(c) {
			l.addTokenIfCan()
			lexNumber(l)
			l.i--
		} else if l.tryOperator() {
		} else {
			l.word.WriteByte(c)
		}
	}

	l.addTokenIfCan()
	close(l.TokChan)
}

func lexNumber(l *Lexer) {
	numType := INTEGER
	numstr := strings.Builder{}

	var digits map[byte]bool = nil
	if c := l.CharAtOffset(0); c == '0' {
		switch nc := l.CharAtOffset(1); nc {
		case 'x':
			{
				l.i += 2
				digits = hexDigits
			}
		case 'b':
			{
				l.i += 2
				digits = binaryDigits
			}
		default:
			{
				// base 8
				// handles corner case of the last char of source code being '0'
				if l.CharAtOffset(1) == 0 {
					digits = decimalDigits
				} else {
					l.i++
					digits = octalDigits
					numstr.WriteByte('0')
				}
			}
		}
	} else {
		// todo: handle expo syntax
		// base 10
		digits = decimalDigits

	}

	if l.CharAtOffset(0) == 0 {
		return
	}
	for ; l.i < len(l.sourceCode); l.i++ {
		c := l.CharAtOffset(0)
		if _, ok := digits[c]; ok {
			numstr.WriteByte(c)
		} else if c == byte(MATH_DOT) && numType == INTEGER {
			numType = FLOATINGPT
			numstr.WriteByte(byte(MATH_DOT))
		} else {
			// erroneous state
			// todo: show error
			break
		}
	}
	if numstr.Len() == 0 {
		// todo: show error
		return
	}
	l.addTokenAndClearWord(NewLexerToken(numType, numstr.String(), l.lineCnt))
}
