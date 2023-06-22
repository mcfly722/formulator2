package explorer

import (
	"encoding/json"
	"math"
)

type Function interface {
	Calculate(childs []*Node) float64
	ToString(childs []*Node) string
}

type Node struct {
	Childs             []*Node
	RecombineFunctions []Function
	result             float64
}

func (node *Node) ToString() string {
	return ""
}

type E struct{}

func (e *E) Calculate(childs []*Node) float64 {
	return math.E
}

func (e *E) ToString(childs []*Node) string {
	return "e"
}

type Pi struct{}

func (pi *Pi) Calculate(childs []*Node) float64 {
	return math.Pi
}

func (pi *Pi) ToString(childs []*Node) string {
	return "pi"
}

type One struct{}

func (one *One) Calculate(childs []*Node) float64 {
	return 1
}

func (one *One) ToString(childs []*Node) string {
	return "1"
}

type Invert struct{}

func (invert *Invert) Calculate(childs []*Node) float64 {
	return -(childs[0].result)
}

func (invert *Invert) ToString(childs []*Node) string {
	return "-" + childs[0].ToString()
}

type Round struct{}

func (round *Round) Calculate(childs []*Node) float64 {
	return math.Round(childs[0].result)
}

func (round *Round) ToString(childs []*Node) string {
	return "Round[" + childs[0].ToString() + "]"
}

type Add struct{}

func (add *Add) Calculate(childs []*Node) float64 {
	return childs[0].result + childs[1].result
}

func (add *Add) ToString(childs []*Node) string {
	return "(" + childs[0].ToString() + "+" + childs[1].ToString() + ")"
}

type Multiply struct{}

func (multiply *Multiply) Calculate(childs []*Node) float64 {
	return childs[0].result * childs[1].result
}

func (multiply *Multiply) ToString(childs []*Node) string {
	return "(" + childs[0].ToString() + "*" + childs[1].ToString() + ")"

}

func (node *Node) Fill() {

	switch len(node.Childs) {
	case 0:
		node.RecombineFunctions = []Function{&E{}, &Pi{}, &One{}}
	case 1:
		node.RecombineFunctions = []Function{&Invert{}, &Round{}}
	case 2:
		node.RecombineFunctions = []Function{&Add{}, &Multiply{}}
	default:
		panic("unsupported number of child nodes for function")
	}

	for _, child := range node.Childs {
		child.Fill()
	}
}

func TreeFromJSON(data []byte) (*Node, error) {
	root := Node{}

	err := json.Unmarshal(data, &root)
	if err != nil {
		return nil, err
	}

	return &root, nil
}

func (node *Node) TreeToJSON() string {
	json, err := json.MarshalIndent(node, "", " ")
	if err != nil {
		panic(err)
	}

	return string(json)
}
