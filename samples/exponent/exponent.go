package main

import (
	"encoding/json"
	"fmt"
	"math"
	"math/rand"
)

// Point ...
type Point struct {
	X float64
	Y float64
}

func main() {

	points := []Point{}

	for i := 0; i < 100; i++ {

		x := rand.Float64()*8 - 4
		y := math.Exp(x)

		points = append(points, Point{X: x, Y: y})
	}

	json, _ := json.MarshalIndent(points, "", "\t")
	fmt.Println(string(json))
}
