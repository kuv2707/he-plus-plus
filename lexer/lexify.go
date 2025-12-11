package lexer

import (
	"bytes"
	"encoding/binary"
	"math"
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
				for ; l.i < len(l.sourceCode) && l.CharAtOffset(0) != '\n'; l.i++ {
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

	var digits map[byte]bool = nil
	var base int64
	if c := l.CharAtOffset(0); c == '0' {
		switch nc := l.CharAtOffset(1); nc {
		case 'x':
			{
				l.i += 2
				digits = hexDigits
				base = 16
			}
		case 'b':
			{
				l.i += 2
				digits = binaryDigits
				base = 2
			}
		default:
			{
				// base 8
				// handles corner case of the last char of source code being '0'
				if l.CharAtOffset(1) == 0 {
					digits = decimalDigits
					base = 10
				} else {
					l.i++
					digits = octalDigits
					base = 8
				}
			}
		}
	} else {
		// todo: handle expo syntax
		// base 10
		digits = decimalDigits
		base = 10

	}

	if l.CharAtOffset(0) == 0 {
		return
	}

	var intPart int64 = 0
	var decPart int64 = 0
	scale := 0
	for ; l.i < len(l.sourceCode); l.i++ {
		c := l.CharAtOffset(0)
		if _, ok := digits[c]; ok {
			if c >= 'A' {
				c -= 'A'
			} else {
				c -= '0'
			}
			if numType == INTEGER {
				intPart = intPart*base + int64(c)
			} else {
				decPart = decPart*int64(base) + int64(c)
				scale++
			}
		} else if c == byte(MATH_DOT) && numType == INTEGER {
			numType = FLOATINGPT

		} else {
			// erroneous state
			// todo: show error
			break
		}
	}
	str := new(bytes.Buffer)
	if numType == INTEGER {
		binary.Write(str, binary.BigEndian, intPart)
	} else {
		var t float64 = float64(intPart) + float64(decPart)/(math.Pow10(scale))
		binary.Write(str, binary.BigEndian, t)
	}
	l.addTokenAndClearWord(NewLexerToken(numType, str.String(), l.lineCnt))
}
