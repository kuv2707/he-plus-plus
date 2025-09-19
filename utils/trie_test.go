package utils_test

import (
	"he++/utils"
	"testing"
)

func TestTrie(t *testing.T) {

	searchFixtures := []string{"abc", "abcde", "bcde", "=", "==", "===", "instance", "instanceof", "and"}
	trie := utils.MakeTrie()
	for i := range searchFixtures {
		trie.Insert(searchFixtures[i])
	}
	for i := range searchFixtures {
		if !trie.Search(searchFixtures[i]) {
			t.Errorf("Insert-Search failed for value %s", searchFixtures[i])
		}
	}
	if trie.Search("abbde") || trie.Search("abcd") {
		t.Errorf("Found non-existant element")
	}

	type Case struct {
		str string
		off int
		exp int
	}
	matchFixtures := []Case{
		{"===", 0, 2},
		{"===", 1, 1},
		{"===", 2, 0},
		{"abcdefghi", 1, 3}, // start from ind 1 match till 4
		{"not", 0, -1},
		{"not", 1, -1},
		{"not", 10, -1},
	}
	for i := range matchFixtures {
		soff := trie.MatchLongest(matchFixtures[i].str, matchFixtures[i].off)
		if soff != matchFixtures[i].exp {
			t.Errorf("Failed MatchLongest for input %s received %d expected %d", matchFixtures[i].str, soff, matchFixtures[i].exp)

		}
	}
}
