package main

import (
	// "fmt"
	"math"

	"github.com/Fito305/tolling/types"
)

// We like to end our interface names with (er).
type CalculatorServicer interface {
	CalculateDistance(types.OBUData) (float64, error)
}

// Implementation of the interface
type CalculatorService struct {
	prevPoint []float64
	// points [][]float64 // a slice of slice // replaced do to endless appending and running out of memory.
}

// Constructor
func NewCalculatorService() CalculatorServicer {
	return &CalculatorService{
		// points: make([][]float64, 0), // refactored due to endless appending of points.
	}
}

func (s *CalculatorService) CalculateDistance(data types.OBUData) (float64, error) {
	distance := 0.0
	if len(s.prevPoint) > 0 {
		distance = calculateDistance(s.prevPoint[0], s.prevPoint[1], data.Lat, data.Long)
	}
	s.prevPoint = []float64{data.Lat, data.Long}
	// if len(s.points) > 0 { // Refactored with a more efficient memory saving implementation.
	// 	prevPoint := s.points[len(s.points) -1] // the last point
	// 	distance = calculateDistance(prevPoint[0] /*lat*/, prevPoint[1] /*long*/, data.Lat, data.Long)
	// }
	// s.points = append(s.points, []float64{data.Lat, data.Long}) // have to fix this. You don't want to keep appending otherwise you will have a memory problem. The RAM will be run out of memory.
	return distance, nil
}

func calculateDistance(x1, y1, x2, y2 float64) float64 {
	return math.Sqrt(math.Pow(x2 - x1, 2) + math.Pow(y2 - y1, 2))
}



// To be a successful programmer you need to know how and where to find things.
