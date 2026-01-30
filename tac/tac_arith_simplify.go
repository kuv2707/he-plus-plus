package tac

import (
	"he++/lexer"
	"math/bits"
)

func (ftac *FunctionTAC) simplifyInstr(tac ThreeAddressInstr) ThreeAddressInstr {
	switch v := tac.(type) {
	case *BinaryOpInstr:
		{
			_, ok1 := v.arg1.(*VRegArg)
			_, ok2 := v.arg2.(*VRegArg)
			if ok1 && ok2 {
				// no scope of simplif
				return v
			} else if ok1 && !ok2 {
				// maybe warn about x/0 in ast validation..
				return simplifyArithmetic(v)
				// return v
			} else if !ok1 && ok2 {
				// num and reg
				// check if the number is 0 and simplify accordingly
				return v
			} else {
				// both numeric, can be precomputed
				instr := &AssignInstr{assnTo: v.assnTo, arg: doArithmetic(v.arg1, v.arg2, v.op)}
				instr.setLabels(tac.Labels())
				return instr
			}
		}
	case *UnaryOpInstr:
		{
			// todo
		}
	}
	return tac
}

func doArithmetic(a, b TACOpArg, op TACOperator) TACOpArg {
	aInt, aIsInt := a.(*ImmIntArg)
	aFloat, _ := a.(*ImmFloatArg)
	bInt, bIsInt := b.(*ImmIntArg)
	bFloat, _ := b.(*ImmFloatArg)

	// Convert to float64 for unified calculation
	// todo: handle separately
	var aVal, bVal float64
	if aIsInt {
		aVal = float64(aInt.num)
	} else {
		aVal = aFloat.num
	}
	if bIsInt {
		bVal = float64(bInt.num)
	} else {
		bVal = bFloat.num
	}

	var result float64
	switch string(op) {
	case lexer.ADD:
		result = aVal + bVal
	case lexer.SUB:
		result = aVal - bVal
	case lexer.MUL:
		result = aVal * bVal
	case lexer.DIV:
		result = aVal / bVal
	}

	// Return as ImmInt if both inputs were int and result is whole number
	if aIsInt && bIsInt && result == float64(int64(result)) {
		return &ImmIntArg{num: int64(result)}
	}
	return &ImmFloatArg{num: result}
}

func simplifyArithmetic(ins *BinaryOpInstr) ThreeAddressInstr {
	// x/2 -> x >> 1 ; x*2 -> x<<1 ; x*0 -> 0
	// todo: x/0 should give error, ideally earlier in the compilation pipeline
	switch v := ins.arg2.(type) {
	case *ImmIntArg:
		if v.num == 0 {
			if string(ins.op) == lexer.MUL {
				return &AssignInstr{assnTo: ins.assnTo, arg: &ImmIntArg{num: 0}}
			} else if string(ins.op) == lexer.ADD || string(ins.op) == lexer.SUB {
				return &AssignInstr{assnTo: ins.assnTo, arg: ins.arg1}
			}
		} else {
			if v.num > 0 && (v.num&(v.num-1)) == 0 {
				// v.num is a power of 2
				if string(ins.op) == lexer.MUL {
					ins.op = TACOperator(lexer.LSHIFT)
					ins.arg2 = &ImmIntArg{num: LogOfTwoPower(v.num)}
				} else if string(ins.op) == lexer.DIV {
					ins.op = TACOperator(lexer.RSHIFT)
					ins.arg2 = &ImmIntArg{num: LogOfTwoPower(v.num)}
				}
			}
		}
	case *ImmFloatArg:
		// todo
	}
	return ins
}

func LogOfTwoPower(num int64) int64 {
	return int64(bits.TrailingZeros(uint(num)))
}
