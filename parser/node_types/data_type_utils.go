package node_types

// nodes for pointer, arr, obj, primitive, error

type DataType interface {
	Text() string
}

type NamedType struct {
	Name string
}

func (dt *NamedType) Text() string {
	return dt.Name
}

type TypePrefix int

const (
	ArrayOf = iota
	PointerOf
)

func (k TypePrefix) String() string {
	switch k {
	case ArrayOf:
		return "[]"
	case PointerOf:
		return "&"
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

type StructType struct {
}

type ErrorType struct {
	Message string
}

func (et *ErrorType) Text() string {
	return "{ERROR: " + et.Message + "}"
}

type VoidType struct {
}

func (vt *VoidType) Text() string {
	return "void"
}
