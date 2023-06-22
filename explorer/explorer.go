package explorer

import (
	zeroOneTwoTree "github.com/formulator2/explorer/step1/zeroOneTwoTree"
)

func GetNextBracketsSequence(brackets string, maxChilds int) (string, error) {
	return zeroOneTwoTree.GetNextBracketsSequence(brackets, maxChilds)
}

func SearchSolution(bracketsSequence string, deviationThreshold float64) (float64, string, error) {

	tree, err := zeroOneTwoTree.BracketsToTree(bracketsSequence)
	if err != nil {
		return 0, "", nil
	}

	json := tree.TreeToJSON()

	node, err := TreeFromJSON(json)
	if err != nil {
		return 0, "", nil
	}

	node.FillFunctions()
	node.FillRecombinationWeights()

	return 0, string(node.TreeToJSON()), nil
}
