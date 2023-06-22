package zeroOneTwoTree

import (
	"encoding/json"
	"fmt"
)

// Node of tree
type Node struct {
	Childs []*Node
}

// BracketStep represent number of opened and closed brackets
type BracketStep struct {
	Opens  int
	Closes int
}

func recursion(_x int, _y int, bracketsStack []BracketStep, maxChilds int, diagonal [32]int, maxBracketPairs int, ready func(bracketsStack []BracketStep, diagonal [32]int)) {
	if _x == maxBracketPairs && _y == maxBracketPairs {
		//ready
		ready(bracketsStack, diagonal)
	} else {
		if diagonal[_x-_y] < maxChilds {
			for x := _x + 1; x <= maxBracketPairs; x++ {
				diagonal[x-_y-1] = diagonal[x-_y-1] + 1
				for y := _y + 1; y <= maxBracketPairs; y++ {
					if y <= x {
						if diagonal[x-y] < maxChilds+1 {
							newBracketsStack := append(bracketsStack, BracketStep{Opens: x - _x, Closes: y - _y})
							recursion(x, y, newBracketsStack, maxChilds, diagonal, maxBracketPairs, ready)
						}
					}
				}
			}
		}
	}
}

// Recombine all binary trees
func Recombine(maxBracketPairs int, maxChilds int, ready func(bracketsStack []BracketStep, diagonal [32]int)) {
	diagonal := [32]int{}
	recursion(0, 0, []BracketStep{}, maxChilds, diagonal, maxBracketPairs, ready)
}

func brackets2points(brackets string) ([]BracketStep, int, error) {
	totalOpens := 0
	totalCloses := 0

	points := []BracketStep{}

	for i := 0; i < len(brackets); {

		opens := 0
		for ; i+opens < len(brackets) && brackets[i+opens] == '('; opens++ {
		}
		i = i + opens
		totalOpens = totalOpens + opens

		if i < len(brackets) && brackets[i] != ')' {
			return nil, 0, fmt.Errorf("unexpected symbol %v <- Expecting '(' or ')'", brackets[:i+1])
		}

		closes := 0
		for ; i+closes < len(brackets) && brackets[i+closes] == ')'; closes++ {
		}
		i = i + closes
		totalCloses = totalCloses + closes

		if i < len(brackets) && brackets[i] != '(' {
			return nil, 0, fmt.Errorf("unexpected symbol %v <- Expecting '(' or ')'", brackets[:i+1])
		}

		points = append(points, BracketStep{Opens: opens, Closes: closes})

		if totalOpens < totalCloses {
			return nil, 0, fmt.Errorf("%v <- total closes=%v are greater than opens=%v", brackets[:i], totalCloses, totalOpens)
		}
	}

	if totalOpens != totalCloses {
		return nil, 0, fmt.Errorf("opened brackets=%v closed brackets=%v should be equal", totalOpens, totalCloses)
	}

	return points, (totalOpens + totalCloses) / 2, nil
}

func recursionNext(srcBracketsStack []BracketStep, dstBracketsStack []BracketStep, _x int, _y int, maxBracketPairs int, maxChilds int, diagonal [32]int, currentRecursionStep int, previousSolutionAlreadyReached bool) ([]BracketStep, bool, bool) {
	if _x == maxBracketPairs && _y == maxBracketPairs {
		if !previousSolutionAlreadyReached && len(srcBracketsStack) > 0 {
			return []BracketStep{}, true, false
		}
		return dstBracketsStack, true, true
	}

	if diagonal[_x-_y] < maxChilds {
		for x := _x + 1; x <= maxBracketPairs; x++ {
			diagonal[x-_y-1] = diagonal[x-_y-1] + 1

			if previousSolutionAlreadyReached || len(srcBracketsStack) == 0 || (x-_x) >= srcBracketsStack[currentRecursionStep].Opens {

				for y := _y + 1; y <= maxBracketPairs; y++ {

					if previousSolutionAlreadyReached || len(srcBracketsStack) == 0 || (y-_y) >= srcBracketsStack[currentRecursionStep].Closes {

						if y <= x {
							if diagonal[x-y] < maxChilds+1 {

								newBracketsStack := append(dstBracketsStack, BracketStep{Opens: x - _x, Closes: y - _y})
								tail, reached, solutionFound := recursionNext(srcBracketsStack, newBracketsStack, x, y, maxBracketPairs, maxChilds, diagonal, currentRecursionStep+1, previousSolutionAlreadyReached)

								previousSolutionAlreadyReached = reached

								if solutionFound {
									return tail, previousSolutionAlreadyReached, solutionFound
								}

							}
						}
					}
				}

			}

		}
	}

	return []BracketStep{}, previousSolutionAlreadyReached, false
}

// GetNextBracketsSequence get current brackets representation of tree and return next one tree in brackets representation
func GetNextBracketsSequence(brackets string, maxChilds int) (string, error) {
	if len(brackets) == 0 {
		return "()", nil
	}

	bracketsStack, maxBracketPairs, err := brackets2points(brackets)
	if err != nil {
		return "", err
	}

	diagonal := [32]int{}

	if len(bracketsStack) == 1 {
		bracketsStack = []BracketStep{}
		maxBracketPairs++
	}

	nextBracketCombination, _, _ := recursionNext(bracketsStack, []BracketStep{}, 0, 0, maxBracketPairs, maxChilds, diagonal, 0, false)

	return BracketsStepsToString(nextBracketCombination), nil
}

// BracketsStepsToString serialize bracketSteps to String
func BracketsStepsToString(tail []BracketStep) string {
	output := ""
	for _, step := range tail {
		for i := 0; i < step.Opens; i++ {
			output += "("
		}
		for i := 0; i < step.Closes; i++ {
			output += ")"
		}
	}
	return output
}

func findMoreThanOneChilds(node *Node) int {
	if len(node.Childs) > 1 {
		return len(node.Childs)
	}
	if len(node.Childs) == 0 {
		return 0
	}
	return findMoreThanOneChilds(node.Childs[0])
}

// BracketsToTree generates expression tree based on string of brackets
func BracketsToTree(input string) (*Node, error) {
	if input == "" {
		return nil, fmt.Errorf("input string could not be empty")
	}

	root, err := bracketsToTree(input)
	if err != nil {
		return nil, err
	}
	if findMoreThanOneChilds(root) < 2 {
		//fmt.Println(fmt.Sprintf("polynom:%v", input))
		return root.Childs[0], nil
	}

	//fmt.Println(fmt.Sprintf("not polynom:%v", input))
	return root, nil
}

func bracketsToTree(input string) (*Node, error) {
	root := Node{Childs: []*Node{}}

	if input == "" {
		return &root, nil
	}

	counter := 0
	from := 0

	for i := 0; i < len(input); i++ {

		if input[i] == '(' {
			if counter == 0 {
				from = i
			}
			counter++
		}

		if input[i] == ')' {
			counter--

			if counter < 0 {
				return nil, fmt.Errorf("%v<- incorrect brackets balance, could not close not opened bracket", input[0:i+1])
			}

			if counter == 0 {
				argument, _ := bracketsToTree(input[from+1 : i])
				if argument != nil {
					root.Childs = append(root.Childs, argument)
				}
			}

		}

		if input[i] != '(' && input[i] != ')' {
			return nil, fmt.Errorf("%v<- unknown symbol", input[0:i+1])
		}
	}

	if counter != 0 {
		return nil, fmt.Errorf("number of opened brackets are not equal to closed (difference=%v)", counter)
	}

	return &root, nil
}

// MaxChilds get maximum childs for tree
func (treeRoot *Node) MaxChilds() int {
	max := len(treeRoot.Childs)

	for _, child := range treeRoot.Childs {
		m := child.MaxChilds()
		if m > max {
			max = m
		}
	}

	return max
}

func (treeRoot *Node) TreeToJSON() []byte {
	json, err := json.MarshalIndent(treeRoot, "", "   ")
	if err != nil {
		panic(err)
	}

	return json
}
