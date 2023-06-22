package explorer

import (
	zeroOneTwoTree "github.com/formulator2/explorer/step1/zeroOneTwoTree"
)

func GetNextBracketsSequence(brackets string, maxChilds int) (string, error) {
	return zeroOneTwoTree.GetNextBracketsSequence(brackets, maxChilds)
}

func SearchSolution(bracketsSequence string, deviationThreshold float64) (float64, string, error) {

	return 0, "", nil
}
