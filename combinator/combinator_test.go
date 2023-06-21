package combinator

import (
	"testing"
)

func Test_Combinator(t *testing.T) {

	sequence, err := GetNextBracketsSequence("()", 2)
	if err != nil {
		t.Fatal(err)
	}

	if sequence != "()()" {
		t.Fatal(sequence)
	}

	t.Logf(sequence)
}
