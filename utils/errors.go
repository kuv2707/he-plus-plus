package utils

import (
	"fmt"
)

type CompilerErrorKind string

const (
	SyntaxError    CompilerErrorKind = "SyntaxError"
	TypeError      CompilerErrorKind = "TypeError"
	UndefinedError CompilerErrorKind = "UndefinedError"
	NotAllowed CompilerErrorKind = "NotAllowed"
)

type CompilerError struct {
	Line int
	Name CompilerErrorKind
	Msg  string
}

func (s *CompilerError) String() string {
	return Red(fmt.Sprintf("%s at line %s: %s", Underline(string(s.Name)), Blue(Underline(fmt.Sprint(s.Line))), s.Msg))
}
