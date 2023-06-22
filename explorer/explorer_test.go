package explorer

import (
	"testing"
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

func Test_SearchSolution1(t *testing.T) {
	deviation, solution, err := SearchSolution("(()())", 10)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("solution %v\n", solution)
	t.Logf("deviation %v\n", deviation)
}
