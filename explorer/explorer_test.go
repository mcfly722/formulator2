package explorer

import (
	"encoding/json"
	"testing"

	zeroOneTwoTree "github.com/formulator2/explorer/step1/zeroOneTwoTree"
)

func Test_GetNextBracketsSequence(t *testing.T) {

	sequence, err := GetNextBracketsSequence("()", 2)
	if err != nil {
		t.Fatal(err)
	}

	if sequence != "()()" {
		t.Fatal(sequence)
	}

	t.Logf(sequence)
}

func Test_BracketsToTree(t *testing.T) {
	sequence := "((())())()"
	root, err := zeroOneTwoTree.BracketsToTree(sequence)
	if err != nil {
		t.Fatal(err)
	}

	json, err := json.MarshalIndent(root, "", " ")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf(string(json))
}

func Test_SearchSolution1(t *testing.T) {
	deviation, solution, err := SearchSolution("(()())", 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("solution %v\n", solution)
	t.Logf("deviation %v\n", deviation)
}
