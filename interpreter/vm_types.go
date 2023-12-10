package interpreter

type Pointer struct {
	address int
	size    int
	scopeId string
	temp bool
}
