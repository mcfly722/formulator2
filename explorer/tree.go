package explorer

import (
	"encoding/json"
	"math"
)

type Function interface {
	Calculate(childs []*Node) float64
	ToString(childs []*Node) string
	MarshalJSON() ([]byte, error)
}

type Node struct {
	Childs             []*Node
	RecombineFunctions []Function
	RecombineWeight    uint64
	result             float64
}

func (node *Node) ToString() string {
	return ""
}

type E struct{}

func (e *E) Calculate(childs []*Node) float64 { return math.E }
func (e *E) ToString(childs []*Node) string   { return "e" }
func (e *E) MarshalJSON() ([]byte, error)     { return json.Marshal("e") }

type Pi struct{}

func (pi *Pi) Calculate(childs []*Node) float64 { return math.Pi }
func (pi *Pi) ToString(childs []*Node) string   { return "pi" }
func (pi *Pi) MarshalJSON() ([]byte, error)     { return json.Marshal("pi") }

type Phi struct{}

func (phi *Phi) Calculate(childs []*Node) float64 { return math.Phi }
func (phi *Phi) ToString(childs []*Node) string   { return "phi" }
func (phi *Phi) MarshalJSON() ([]byte, error)     { return json.Marshal("phi") }

type One struct{}

func (one *One) Calculate(childs []*Node) float64 { return 1 }
func (one *One) ToString(childs []*Node) string   { return "1" }
func (one *One) MarshalJSON() ([]byte, error)     { return json.Marshal("1") }

type Invert struct{}

func (invert *Invert) Calculate(childs []*Node) float64 { return -(childs[0].result) }
func (invert *Invert) ToString(childs []*Node) string   { return "-" + childs[0].ToString() }
func (invert *Invert) MarshalJSON() ([]byte, error)     { return json.Marshal("invert") }

type Round struct{}

func (round *Round) Calculate(childs []*Node) float64 { return math.Round(childs[0].result) }
func (round *Round) ToString(childs []*Node) string   { return "Round[" + childs[0].ToString() + "]" }
func (round *Round) MarshalJSON() ([]byte, error)     { return json.Marshal("round") }

type Add struct{}

func (add *Add) Calculate(childs []*Node) float64 { return childs[0].result + childs[1].result }
func (add *Add) ToString(childs []*Node) string {
	return "(" + childs[0].ToString() + "+" + childs[1].ToString() + ")"
}
func (add *Add) MarshalJSON() ([]byte, error) { return json.Marshal("+") }

type Multiply struct{}

func (multiply *Multiply) Calculate(childs []*Node) float64 {
	return childs[0].result * childs[1].result
}

func (multiply *Multiply) ToString(childs []*Node) string {
	return "(" + childs[0].ToString() + "*" + childs[1].ToString() + ")"
}
func (multiply *Multiply) MarshalJSON() ([]byte, error) { return json.Marshal("*") }

func (node *Node) FillFunctions() uint64 {

	switch len(node.Childs) {
	case 0:
		node.RecombineFunctions = []Function{&E{}, &Pi{}, &Phi{}, &One{}}
	case 1:
		node.RecombineFunctions = []Function{&Invert{}, &Round{}}
	case 2:
		node.RecombineFunctions = []Function{&Add{}, &Multiply{}}
	default:
		panic("unsupported number of child nodes for function")
	}

	summaryWeigth := (uint64)(len(node.RecombineFunctions))

	for _, child := range node.Childs {
		weigth := child.FillFunctions()
		summaryWeigth *= weigth
	}

	node.RecombineWeight = summaryWeigth
	return node.RecombineWeight
}

func (node *Node) FillRecombinationWeights() uint64 {

	summaryWeigth := (uint64)(len(node.RecombineFunctions))

	for _, child := range node.Childs {
		weigth := child.FillRecombinationWeights()
		summaryWeigth *= weigth
	}

	node.RecombineWeight = summaryWeigth
	return node.RecombineWeight
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

func (node *Node) GetAllSubFunctions0() []*Node {
	if len(node.Childs) == 0 {
		return []*Node{node}
	}

	childs := []*Node{}
	for _, child := range node.Childs {
		childs = append(childs, child.GetAllSubFunctions0()...)
	}

	return childs
}
