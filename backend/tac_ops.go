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

type TACOpArg interface {
	LocType() TACLocationType
	String() string
}

type VRegArg struct {
	regNo VirtualRegisterNumber
}

func (k *VRegArg) LocType() TACLocationType {
	return VReg
}

func (arg *VRegArg) String() string {
	return fmt.Sprintf("%s%d", utils.Yellow(VReg.String()), arg.regNo)
}

type ImmIntArg struct {
	num int64
	// size     int
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

func (k *ImmFloatArg) LocType() TACLocationType {
	return Imm
}

func (arg *ImmFloatArg) String() string {
	return fmt.Sprintf("%s%d", utils.Yellow("#"), arg.num)
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
