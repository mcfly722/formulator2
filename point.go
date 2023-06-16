package main

import (
	"encoding/json"
	"os"
)

type Point struct {
	X float64
	Y float64
}

func loadPointsFromFile(samplesFile string) ([]Point, error) {
	body, err := os.ReadFile(samplesFile)
	if err != nil {
		return nil, err
	}

	points := []Point{}

	err = json.Unmarshal(body, &points)
	if err != nil {
		return nil, err
	}

	return points, nil
}
