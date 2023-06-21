package combinator

import (
	zeroOneTwoTree "github.com/formulator2/combinator/step1/zeroOneTwoTree"
)

func GetNextBracketsSequence(brackets string, maxChilds int) (string, error) {
	return zeroOneTwoTree.GetNextBracketsSequence(brackets, maxChilds)
}
