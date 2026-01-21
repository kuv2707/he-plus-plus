package tac

import (
	"fmt"
	"he++/parser/node_types"
	staticanalyzer "he++/static_analyzer"
	"he++/utils"
)

type TACLocationType byte
type VirtualRegisterNumber = int64
type TACOperator string

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

func (dc DataCategory) IsFloating() bool {
	return dc == F32 || dc == F64
}

func (dc DataCategory) SizeBytes() int {
	switch dc {
	case I16:
		return 2
	case I32, F32:
		return 4
	case I64, F64, PTR, AGGREGATE:
		return 8
	case BYTE:
		return 1
	case VOID:
		return 0
	default:
		return 0
	}
}

type TACOpArg interface {
	LocType() TACLocationType
	String() string
}

type VRegArg struct {
	RegNo VirtualRegisterNumber
}

func (k *VRegArg) LocType() TACLocationType {
	return VReg
}

func (arg *VRegArg) String() string {
	return fmt.Sprintf("%s%d", utils.Yellow(VReg.String()), arg.RegNo)
}

type ImmIntArg struct {
	num int64
	// size     int
}

func (k ImmIntArg) Num() int64 {
	return k.num
}

func (k *ImmIntArg) LocType() TACLocationType {
	return Imm
}

func (arg *ImmIntArg) String() string {
	return fmt.Sprintf("%s%d", utils.Yellow("#"), arg.num)
}

type ImmFloatArg struct {
	num float64
	// size     int
}

func (k ImmFloatArg) Num() float64 {
	return k.num
}

func (k *ImmFloatArg) LocType() TACLocationType {
	return Imm
}

func (arg *ImmFloatArg) String() string {
	return fmt.Sprintf("%s%f", utils.Yellow("#"), arg.num)
}

type NULLOpArg struct{}

func (k *NULLOpArg) LocType() TACLocationType {
	return Null
}

func (arg *NULLOpArg) String() string {
	return utils.Yellow("_")
}

var NOWHERE TACOpArg = &NULLOpArg{}

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
			return I32
		} else {
			return F32
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
