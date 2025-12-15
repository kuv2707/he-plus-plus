package backend

import (
	"fmt"
	"he++/parser/node_types"
	staticanalyzer "he++/static_analyzer"
	"he++/utils"
)

type TACLocationType byte
type VirtualRegisterNumber = int64

const (
	VReg TACLocationType = iota
	Imm
	Null
)

type DataCategory int16

const (
	I16 DataCategory = iota
	I32
	I64
	F32
	F64
	PTR
	BYTE
	AGGREGATE // maybe functionally eq to PTR
	VOID
)

type TACOpArg struct {
	locType TACLocationType
	ival    int64
}

func (arg TACOpArg) String() string {
	return fmt.Sprintf("%s%d", utils.Yellow(arg.locType.String()), arg.ival)
}

var NOWHERE = TACOpArg{Null, 0}

func (t TACLocationType) String() string {
	switch t {
	case VReg:
		return "R"
	case Imm:
		return "#"
	case Null:
		return "XX"
	}
	panic("Undefined TACLocationType")
}

func dataCategoryForType(dt node_types.DataType) DataCategory {
	switch v := dt.(type) {
	case *node_types.PrefixOfType:
		{
			if v.Prefix == node_types.ArrayOf {
				return PTR
			} else if v.Prefix == node_types.PointerOf {
				return PTR
			}
		}
	case *node_types.StructType:
		return PTR
	}

	switch dt.Size() {
	case 1:
		return BYTE
	case 2:
		return I16
	case 4:
		if dt.TypeId() == staticanalyzer.INT_DATATYPE.TypeId() {
			return F32
		} else {
			return I32
		}
	case 8:
		if dt.TypeId() == staticanalyzer.FLOAT_DATATYPE.TypeId() {
			return F64
		} else {
			return I64
		}
	default:
		return AGGREGATE
	}
}
