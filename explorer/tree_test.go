package explorer

import (
	"encoding/json"
	"fmt"
	"testing"

	zeroOneTwoTree "github.com/formulator2/explorer/step1/zeroOneTwoTree"
)

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

func Test_GetAllSubFunctions0(t *testing.T) {
	tree, err := zeroOneTwoTree.BracketsToTree("(((()))())(()())")
	if err != nil {
		t.Fatal(err.Error())
	}

	js1 := tree.TreeToJSON()

	node, err := TreeFromJSON(js1)
	if err != nil {
		t.Fatal(err.Error())
	}

	node.FillFunctions()
	node.FillRecombinationWeights()

	//t.Logf(string(node.TreeToJSON()))

	functions0 := node.GetAllSubFunctions0()

	js2, err := json.MarshalIndent(functions0, "", " ")
	if err != nil {
		panic(err)
	}

	t.Logf(string(js2))
}

func Test_RecombineNodes(t *testing.T) {
	nodes := []*Node{
		{result: 5},
		{result: 5},
		{result: 5},
		{result: 5},
		{result: 5},
	}

	//fmt.Printf("%v\n", nodes)

	print := func(occupied []*Node, free []*Node) {

		for _, node := range free {
			(*node).result = 0
		}

		for _, node := range occupied {
			(*node).result = 1
		}

		for _, node := range nodes {
			fmt.Printf("%v", (*node).result)
		}

		fmt.Printf("\n")
	}

	RecombineNodes(&nodes, 0, 2, print)
}
