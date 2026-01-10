package tac

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
				// todo: impl x/2 -> x >> 1 and x*2 -> x<<1, x*0 -> 0
				// maybe warn about x/0 in ast validation..
				return v
			} else if !ok1 && ok2 {
				// num and reg
				// check if the number is 0 and simplify accordingly
				return v
			} else {
				// both numeric, can be precomputed
				return &AssignInstr{assnTo: v.assnTo, arg: doArithmetic(v.arg1, v.arg2, v.op)}
			}
		}
	case *UnaryOpInstr:
		{
			// todo
		}
	}
	return tac
}

func doArithmetic(a TACOpArg, b TACOpArg, op string) TACOpArg {
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
	switch op {
	case "+":
		result = aVal + bVal
	case "-":
		result = aVal - bVal
	case "*":
		result = aVal * bVal
	case "/":
		result = aVal / bVal
	}

	// Return as ImmInt if both inputs were int and result is whole number
	if aIsInt && bIsInt && result == float64(int64(result)) {
		return &ImmIntArg{num: int64(result)}
	}
	return &ImmFloatArg{num: result}
}
