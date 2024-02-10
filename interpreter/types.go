package interpreter

import "he++/parser"

type DataType byte

const (
	NULL DataType = iota
	NUMBER 
	CHAR
	BOOLEAN
	POINTER
	STRING
	ARRAY
	STRUCT
)

var typeNameMap = []string{"NUMBER", "CHAR", "BOOLEAN", "POINTER", "STRING", "ARRAY", "STRUCT"}


func (dt DataType) String() string {
	return typeNameMap[dt]
}

type Reason string

type ScopeContext struct {
	scopeId     string
	scopeTyp    string
	scopeName   string
	variables   map[string]*Pointer
	functions   map[string]parser.TreeNode
	returnValue *Pointer
}
