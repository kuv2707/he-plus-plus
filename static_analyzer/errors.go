package staticanalyzer

import (
	"fmt"
	"he++/utils"
)

type ErrorName string

const (
	SyntaxError    ErrorName = "SyntaxError"
	TypeError      ErrorName = "TypeError"
	UndefinedError ErrorName = "UndefinedError"
)

type StaticAnalyserError struct {
	line int
	name ErrorName
	msg  string
}

func (s *StaticAnalyserError) DisplayString() string {
	return utils.Red(fmt.Sprintf("%s at line %d: %s", s.name, s.line, s.msg))
}
