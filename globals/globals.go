package globals


//todo bitwise operators

var Booleans = []string{"true", "false"}

//contains mapping of number strings to number values and strings to their hash values
var NumMap = map[string][]byte{}

const MOD = 1000000007
const PRIME = 51

func HashString(s string) int {
	var hash int = 0
	for i := 0; i < len(s); i++ {
		hash = (hash*PRIME + int(s[i])) % MOD
	}
	return hash
}

func BeginsWithCapital(s string) bool {
	return s[0] >= 'A' && s[0] <= 'Z'
}
