package utils

import "strconv"

func IsNumber(temp string) bool {
	for i := 0; i < len(temp); i++ {
		if temp[i] < '0' || temp[i] > '9' {
			return false
		}
	}
	return true
}

func StringToInt(str string) int {
	val,err:= strconv.Atoi(str)
	if err!=nil{
		panic(err)
	}
	return val
}

func IsOneOf(temp string,options string) bool {
	for i := 0; i < len(options); i++ {
		if temp == string(options[i]) {
			return true
		}
	}
	return false
}


func IsOneOfArr(str string, options []string)bool{
	for _,word:=range options{
		if( word==str){

			return true
		}
	}
	return false
}
var QUOTES="`"
//todo use regex
func InQuotes(s string)bool{
	return IsOneOf(s[0:1],QUOTES) && IsOneOf(s[len(s)-1:len(s)],QUOTES)
}