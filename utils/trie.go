package utils

// minimal implementation only suitable for inserting elements
type Trie struct {
	Elems []TrieElement
}

type TrieElement struct {
	Stop   bool
	Childs map[byte]int
}

func MakeTrie() Trie {
	return Trie{Elems: []TrieElement{{Childs: make(map[byte]int)}}}
}

func (t *Trie) Insert(s string) {
	root := &t.Elems[0]
	for i := 0; i < len(s); i++ {
		idx := root.Childs[s[i]]
		if idx != 0 {
			root = &t.Elems[idx]
		} else {
			t.Elems = append(t.Elems, TrieElement{Childs: make(map[byte]int)})
			root.Childs[s[i]] = len(t.Elems) - 1
			root = &t.Elems[len(t.Elems)-1]
		}
	}
	root.Stop = true
}

func (t *Trie) Search(s string) bool {
	root := t.Elems[0]
	for i := 0; i < len(s); i++ {
		idx := root.Childs[s[i]]
		if idx != 0 {
			root = t.Elems[idx]
		} else {
			return false
		}
	}
	return root.Stop
}

func (t *Trie) MatchLongest(s string, offset int) int {
	lastI := -1
	root := t.Elems[0]
	for i := 0; i+offset < len(s); i++ {
		idx := root.Childs[s[i+offset]]
		if idx != 0 {
			root = t.Elems[idx]
			if root.Stop {
				lastI = i
			}
		} else {
			break
		}
	}
	return lastI
}
