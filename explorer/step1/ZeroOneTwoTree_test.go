package zeroOneTwoTree

import (
	"encoding/json"
	"testing"
)

func Test_Recombines(t *testing.T) {
	i := 1

	ready := func(bracketsStack []BracketStep, diagonal [32]int) {
		representation := ""
		for _, point := range bracketsStack {
			for opens := 0; opens < point.Opens; opens++ {
				representation += "("
			}
			for closes := 0; closes < point.Closes; closes++ {
				representation += ")"
			}
		}

		t.Logf("%5v   %v   %v", i, representation, diagonal)
		i++
	}

	Recombine(4, 2, ready)
}

func testBracketsForError(t *testing.T, brackets string) {
	_, err := GetNextBracketsSequence(brackets, 2)
	if err != nil {
		t.Logf("correct error handling for %v -> %v", brackets, err)
	} else {
		t.Errorf("GetNextTree('%v') not returned error", brackets)
	}
}

func Test_BracketsOpensCloses(t *testing.T) {
	testBracketsForError(t, "(())(")
}

func Test_BracketsUnexpectedSymbol1(t *testing.T) {
	testBracketsForError(t, "(())!()")
}

func Test_BracketsUnexpectedSymbol2(t *testing.T) {
	testBracketsForError(t, "(())(!)")
}

func Test_BracketsUnexpectedSymbol3(t *testing.T) {
	testBracketsForError(t, "(())()!")
}

func Test_BracketsUnexpectedSymbol4(t *testing.T) {
	testBracketsForError(t, "!(())()")
}

func Test_BracketsClosesGreaterThanOpens(t *testing.T) {
	testBracketsForError(t, "((())))()")
}

func Test_BracketsToTree(t *testing.T) {
	bracketSequence := "()((())())"
	expression, err := BracketsToTree(bracketSequence)
	if err != nil {
		t.Errorf("Cant build expression tree for %v. Error: %v", bracketSequence, err)
	}

	bytes, err := json.Marshal(expression)
	if err != nil {
		t.Errorf("Can't serialize %v. Error:%v", bracketSequence, err)
	}
	t.Log(string(bytes))

}

func testBracketsToTreeError(t *testing.T, testSequence string) {
	_, err := BracketsToTree(testSequence)

	if err == nil {
		t.Errorf("error for '%v' sequence does not catched", testSequence)
	} else {
		t.Logf("error for '%v' successfully catched.\nerror description:%v", testSequence, err)
	}
}

func testBracketsToTreeSuccess(t *testing.T, testSequence string) {
	expression, err := BracketsToTree(testSequence)

	if err != nil {
		t.Errorf("error for %v :%v", testSequence, err)
	}

	bytes, err := json.Marshal(expression)
	if err != nil {
		t.Errorf("Can't serialize %v. Error:%v", testSequence, err)
	}
	t.Log(string(bytes))
}

func Test_BracketsToTree_IncorrectSymbolFromStart(t *testing.T) {
	testBracketsToTreeError(t, "a(()(()))")
}

func Test_BracketsToTree_IncorrectSymbolAtTheEnd(t *testing.T) {
	testBracketsToTreeError(t, "(()(()))b")
}

func Test_BracketsToTree_IncorrectSymbolInTheMiddle(t *testing.T) {
	testBracketsToTreeError(t, "(()((c)))")
}

func Test_BracketsToTree_LostOpeningBracket(t *testing.T) {
	testBracketsToTreeError(t, "())")
}

func Test_BracketsToTree_LostClosingBracket(t *testing.T) {
	testBracketsToTreeError(t, "(()")
}

func Test_BracketsToTree_WrongBracketsSequence(t *testing.T) {
	testBracketsToTreeError(t, "())(")
}

func Test_BracketsToTree_FirstBracket(t *testing.T) {
	testBracketsToTreeSuccess(t, "()")
}

func iterateOverTrees(t *testing.T, bracketSequence string, n int, maxChilds int) {
	for i := 1; i < n; i++ {

		tree, err := BracketsToTree(bracketSequence)
		if err != nil {
			t.Errorf("BracketsToTree ('%v') returned error: %v", bracketSequence, err)
		}

		max := tree.MaxChilds()

		nextBracketSequcence, err := GetNextBracketsSequence(bracketSequence, maxChilds)
		if err != nil {
			t.Errorf("GetNextBracketsSequence ('%v') returned error: %v", bracketSequence, err)
		}

		t.Logf("%3v) %3v %v -> %v", i, max, bracketSequence, nextBracketSequcence)
		bracketSequence = nextBracketSequcence
	}
}

func Test_IterateOverPossibleTrees(t *testing.T) {
	iterateOverTrees(t, "()", 60, 10000)
}

func Test_IterateOverZeroOneTwoTrees(t *testing.T) {
	iterateOverTrees(t, "()", 60, 2)
}

func testBracketsToTree(t *testing.T, bracketsSequence string) {
	astTree, err := BracketsToTree(bracketsSequence)
	if err != nil {
		t.Errorf("BracketsToTree ('%v') returned error: %v", bracketsSequence, err)
	}

	astJSON, _ := json.MarshalIndent(astTree, "", "  ")
	t.Logf("sequence:%v\n%v", bracketsSequence, string(astJSON))

}

func Test_BracketsToTree1(t *testing.T) {
	testBracketsToTree(t, "()")
}
func Test_BracketsToTree2(t *testing.T) {
	testBracketsToTree(t, "(())")
}
func Test_BracketsToTree3(t *testing.T) {
	testBracketsToTree(t, "((()))")
}
func Test_BracketsToTree4(t *testing.T) {
	testBracketsToTree(t, "()()")
}
func Test_BracketsToTree5(t *testing.T) {
	testBracketsToTree(t, "()()()")
}
func Test_BracketsToTree6(t *testing.T) {
	testBracketsToTree(t, "(())()")
}
func Test_BracketsToTree7(t *testing.T) {
	testBracketsToTree(t, "(())(())")
}

func Test_TreeToJSON(t *testing.T) {
	root, err := BracketsToTree("()((())())((()((())())))")
	if err != nil {
		t.Errorf(err.Error())
	}

	t.Logf(string(root.TreeToJSON()))
}

// go test -bench .

func BenchmarkGetNextBracketsSequenceForTwo(b *testing.B) {
	currentSequence := "()"
	for n := 0; n < b.N; n++ {
		next, _ := GetNextBracketsSequence(currentSequence, 2)
		currentSequence = next
	}
}

func BenchmarkBracketsToTree(b *testing.B) {
	currentSequence := "()((())())((()((())())))"
	for n := 0; n < b.N; n++ {
		_, err := BracketsToTree(currentSequence)
		if err != nil {
			b.Errorf("error:%v", err)
		}
	}
}
