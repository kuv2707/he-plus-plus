package interpreter

import "he++/parser"

type DataType byte

const (
	NUMBER DataType = iota
	CHAR
	BOOLEAN
	POINTER
	STRING
	ARRAY
	STRUCT
)

func (dt DataType) String() string {
	return [...]string{"NUMBER", "CHAR", "BOOLEAN", "POINTER", "STRING", "ARRAY", "STRUCT"}[dt]
}

type Reason string

type scopeContext struct {
	scopeId     string
	scopeTyp    string
	scopeName   string
	variables   map[string]*Pointer
	functions   map[string]parser.TreeNode
	returnValue *Pointer
}
