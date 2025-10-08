package staticanalyzer

import (
	"he++/lexer"
	"he++/parser/node_types"
)

func isIndexable(a *Analyzer, dataT node_types.DataType, indexerT node_types.DataType) (node_types.DataType, bool) {
	arrT, ok := dataT.(*node_types.PrefixOfType)
	if !ok {
		return nil, false
	}
	// todo: include string
	if arrT.Prefix != node_types.ArrayOf {
		return nil, false
	}
	// todo: a better way of deducing numeric type
	// a.definedTypes[]
	if indexerT.Text() != lexer.INT {
		return nil, false
	}
	return arrT.OfType, true

}
