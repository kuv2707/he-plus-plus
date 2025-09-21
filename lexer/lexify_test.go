package lexer

import (
	"testing"
)

func TestLexify(t *testing.T) {
	t.Run("Basic assignment", func(t *testing.T) {
		testLexerExpectTokens(t, "definir x = 10", []LexerToken{
			{"keyword", "definir", 1},
			{"identifier", "x", 1},
			{"operator", "=", 1},
			{"int", "10", 1},
		})
	})

	t.Run("Array in object declaration", func(t *testing.T) {
		testLexerExpectTokens(t, "definir a = { name: \" Kislay \" ,\n roll: [ 12, 13 ] };", []LexerToken{
			{"keyword", "definir", 1},
			{"identifier", "a", 1},
			{"operator", "=", 1},
			{"bracket", "{", 1},
			{"identifier", "name", 1},
			{"punctuation", ":", 1},
			{"string_literal", " Kislay ", 1},
			{"punctuation", ",", 1},
			{"identifier", "roll", 2},
			{"punctuation", ":", 2},
			{"bracket", "[", 2},
			{"int", "12", 2},
			{"punctuation", ",", 2},
			{"int", "13", 2},
			{"bracket", "]", 2},
			{"bracket", "}", 2},
			{"punctuation", ";", 2},
		})
	})

	t.Run("Expressions: Unary operators", func(t *testing.T) {
		testLexerExpectTokens(t, "a+++b--*c/d%e", []LexerToken{
			{"identifier", "a", 1},
			{"operator", "++", 1},
			{"operator", "+", 1},
			{"identifier", "b", 1},
			{"operator", "--", 1},
			{"operator", "*", 1},
			{"identifier", "c", 1},
			{"operator", "/", 1},
			{"identifier", "d", 1},
			{"operator", "%", 1},
			{"identifier", "e", 1},
		})
	})

	t.Run("Expressions: Binary operators", func(t *testing.T) {
		testLexerExpectTokens(t, "a+3*b(4,5--) == 0", []LexerToken{
			{"identifier", "a", 1},
			{"operator", "+", 1},
			{"int", "3", 1},
			{"operator", "*", 1},
			{"identifier", "b", 1},
			{"bracket", "(", 1},
			{"int", "4", 1},
			{"punctuation", ",", 1},
			{"int", "5", 1},
			{"operator", "--", 1},
			{"bracket", ")", 1},
			{"operator", "==", 1},
			{"int", "0", 1},
		})
	})
	t.Run("String with escape sequences", func(t *testing.T) {
		testLexerExpectTokens(t, "\"Hello,\\f \\nWor\\tld!\\z\"", []LexerToken{
			{"string_literal", "Hello,\f \nWor\tld!", 1},
		})
	})

	t.Run("Floating point numbers", func(t *testing.T) {
		testLexerExpectTokens(t, "definir x = 10.5", []LexerToken{
			{"keyword", "definir", 1},
			{"identifier", "x", 1},
			{"operator", "=", 1},
			{"floatingpt", "10.5", 1},
		})

		testLexerExpectTokens(t, "definir x = 5.56.78", []LexerToken{
			{"keyword", "definir", 1},
			{"identifier", "x", 1},
			{"operator", "=", 1},
			{"floatingpt", "5.56", 1},
			{"operator", ".", 1},
			{"int", "78", 1},
		})

	})

	t.Run("Numbers in different bases", func(t *testing.T) {
		testLexerExpectTokens(t, "0xDEADBEEF 0123 012389 0xA.5.076  0xYEAH 099 009DEH", []LexerToken{
			{"int", "DEADBEEF", 1},
			{"int", "0123", 1},
			{"int", "0123", 1},
			{"int", "89", 1},
			{"floatingpt", "A.5", 1},
			{"operator", ".", 1},
			{"int", "076", 1},
			{"identifier", "YEAH", 1},
			{"int", "0", 1},
			{"int", "99", 1},
			{"int", "00", 1},
			{"int", "9", 1},
			{"identifier", "DEH", 1},
		})

		testLexerExpectTokens(t, " 0xYEAH 099 0b301", []LexerToken{
			{"identifier", "YEAH", 1},
			{"int", "0", 1},
			{"int", "99", 1},
			{"int", "301", 1},
		})

	})

	// todo: Add tests for warnings
	if !t.Failed() {
		t.Log("\033[32mAll tests passed\033[0m")
	}
}

func testLexerExpectTokens(t *testing.T, sourceCode string, expectedTokens []LexerToken) {
	lexer := LexerOf(sourceCode)
	tokens := lexer.Lexify()

	if len(tokens) != len(expectedTokens) {
		t.Log("\033[31mFailed!\033[0m Tokens:", tokens)
		t.Fatalf("expected %d tokens, got %d", len(expectedTokens), len(tokens))
	}

	for i, token := range tokens {
		if token != expectedTokens[i] {
			t.Log("\033[31mFailed!\033[0m Tokens:", tokens)
			t.Errorf("expected token %v, got %v", expectedTokens[i], token)
		}
	}
}
