package parser
import (
	"toylingo/lexer"
	"toylingo/utils"
	"fmt"
)

var i, length int = 0, 0

func next() lexer.TokenType {
	i++
	return tokensArr[i-1]
}
func prev() lexer.TokenType {
	i--
	return tokensArr[i+1]
}

func peek() lexer.TokenType {
	//todo check index before accessing
	return tokensArr[i+1]
	
}

func index() int {
	return i
}

func maxIndex() int {
	return length
}

func expect(tokenType string) {
	if tokensArr[i].Type != tokenType {
		panic("unexpected token" + tokensArr[i].Type + " " + tokensArr[i].Ref + " expected " + tokenType)
	}
}

// collects tokens till a token of type tokenType is found and consumes but does not include it in returned array
func collectTill(tokenType string) []lexer.TokenType {
	tokens := make([]lexer.TokenType, 0)
	for ; i < len(tokensArr); i++ {
		if tokensArr[i].Type == tokenType {
			break
		}
		tokens = append(tokens, tokensArr[i])
	}
	i++
	return tokens
}

func matchCurrent(tokenType string) bool {
	if index() >= maxIndex() {
		return false
	}
	if tokensArr[i].Type == tokenType {
		return true
	}
	return false
}

func consume(tokenType string) {
	if index() >= maxIndex() {
		return
	}
	if tokensArr[i].Type == tokenType {
		i++
	} else {
		panic("unexpected token" + tokensArr[i].Type + " " + tokensArr[i].Ref + " expected " + tokenType)
	}
}



func (treeNode *TreeNode) PrintTree(space string) {
	color := utils.GetRandomColor()
	fmt.Print(color)
	fmt.Println(space + "{")
	space += "  "
	fmt.Println(space + "desc:" + treeNode.Description)
	fmt.Println(space + "label:" + treeNode.Label)
	for key, val := range treeNode.Properties {
		fmt.Println(space + key + ":")
		val.PrintTree(space + utils.ONETAB)
	}
	if len(treeNode.Children)>0{
		fmt.Println(space + "children:\n"+space+"[")
		for _, child := range treeNode.Children {
	
			child.PrintTree(space + utils.ONETAB)
			fmt.Print(utils.Colors["RESET"])
		}

	}
	fmt.Print(color)
	fmt.Println(space + "]\n" + space[0:len(space)-2] + "}")

}
func printTokensArr(tokens []lexer.TokenType) {
	for i := 0; i < len(tokens); i++ {
		fmt.Print(tokens[i].Ref, " ")
	}
	fmt.Println()
}