package lexer

func (l *Lexer) Lexify() []LexerToken {
	for ; l.i < len(l.sourceCode); l.i++ {
		c := l.sourceCode[l.i]
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
				if l.sourceCode[l.i] == '\\' {
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
			if l.i+1 < len(l.sourceCode) && l.sourceCode[l.i+1] == '/' {
				for ; l.i < len(l.sourceCode) && l.sourceCode[l.i] != '\n'; l.i++ {
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
			// todo: handle hex, octal, expo syntax
			for ; l.i < len(l.sourceCode) && isNumberPart(l.sourceCode[l.i]); l.i++ {
				l.word.WriteByte(l.sourceCode[l.i])
			}
			l.addTokenAndClearWord(NewLexerToken(INTEGER, l.word.String(), l.lineCnt))
			l.i--

			// improve this
		} else if l.tryOperator() {
		} else {
			l.word.WriteByte(c)
		}
	}

	l.addTokenIfCan()
	return l.tokens
}
