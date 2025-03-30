package node_types

import "he++/lexer"

// nodes for pointer, arr, obj, primitive, error

type DataType interface {
	Text() string
	// NumBytes() int
	Equals(DataType) bool
}

type NamedType struct {
	Name string
}

func (dt *NamedType) Text() string {
	return dt.Name
}

func (dt *NamedType) Equals(other DataType) bool {
	if ont, ok := other.(*NamedType); ok {
		return dt.Name == ont.Name
	}
	return false
}

type TypePrefix int

const (
	ArrayOf = iota
	PointerOf
	Dereference
	Unknown
)

func (k TypePrefix) String() string {
	switch k {
	case ArrayOf:
		return "[]"
	case PointerOf:
		return lexer.AMP
	case Dereference:
		// should distinguish between mul and deref ops at lexer level
		return lexer.MUL
	default:
		return "Unknown"
	}
}

type PrefixOfType struct {
	Prefix TypePrefix
	OfType DataType
}

func (dt *PrefixOfType) Text() string {
	return dt.Prefix.String() + dt.OfType.Text()
}

func (dt *PrefixOfType) Equals(other DataType) bool {
	if ont, ok := other.(*PrefixOfType); ok {
		return dt.Prefix.String() == ont.Prefix.String() && dt.OfType.Equals(ont.OfType)
	}
	return false
}

type FuncType struct {
	ReturnType DataType
	ArgTypes   []DataType
}

func (ft *FuncType) Equals(dt DataType) bool {
	oft, ok := dt.(*FuncType)
	if !ok {
		return false
	}
	if len(oft.ArgTypes) != len(ft.ArgTypes) {
		return false
	}
	for i := range ft.ArgTypes {
		if !ft.ArgTypes[i].Equals(oft.ArgTypes[i]) {
			return false
		}
	}
	return true
}

func (ft *FuncType) Text() string {
	ans := "func("
	for _, t := range ft.ArgTypes {
		ans += t.Text() + ","
	}
	ans += ") "
	ans += ft.ReturnType.Text()
	return ans
}

type StructType struct {
	// todo
}

type ErrorType struct {
	Message string
}

func (et *ErrorType) Text() string {
	return "{ERROR: " + et.Message + "}"
}

func (dt *ErrorType) Equals(other DataType) bool {
	// error types are never equal
	return false
}

type VoidType struct {
}

func (vt *VoidType) Text() string {
	return "void"
}

func (dt *VoidType) Equals(other DataType) bool {
	_, ok := other.(*VoidType)
	return ok
}
