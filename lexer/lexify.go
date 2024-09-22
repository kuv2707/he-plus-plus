package lexer


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
				l.lineCnt++
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
				l.popAndInsert(2, NewLexerToken(FLOATINGPT, l.lookBack(2).Text()+DOT+l.word, l.lineCnt))
				l.clearWord()
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
