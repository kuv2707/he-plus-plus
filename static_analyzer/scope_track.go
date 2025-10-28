package staticanalyzer

type ScopeType int

const (
	FUNCTION ScopeType = iota
	CONDITIONAL
	LOOP
	NESTED
	BASE
)

type ScopeEntry struct {
	ScopeType ScopeType
	DefinedTypes map[string]bool
	DefinedSyms map[string]bool
}