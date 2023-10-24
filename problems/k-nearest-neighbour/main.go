package main

import (
	"fmt"
	"sort"
)

type Point struct {
	X float64
	Y float64
}

type PointsByDist struct {
	Points         []Point
	ReferencePoint Point
}

func (p PointsByDist) Len() int {
	return len(p.Points)
}

func (p PointsByDist) Less(i, j int) bool {
	d1 := (p.Points[i].X-p.ReferencePoint.X)*(p.Points[i].X-p.ReferencePoint.X) + (p.Points[i].Y-p.ReferencePoint.Y)*(p.Points[i].Y-p.ReferencePoint.Y)
	d2 := (p.Points[j].X-p.ReferencePoint.X)*(p.Points[j].X-p.ReferencePoint.X) + (p.Points[j].Y-p.ReferencePoint.Y)*(p.Points[j].Y-p.ReferencePoint.Y)

	return d1 < d2
}

func (p PointsByDist) Swap(i, j int) {
	p.Points[i], p.Points[j] = p.Points[j], p.Points[i]
}

func main() {
	startingPoint := Point{3.0, 4.0}

	candidatePoints := []Point{
		{12.0, 18.0},
		{2.0, 5.0},
		{6.0, 9.0},
		{21.0, 14.0},
		{15.0, 19.0},
	}

	k := 3

	pd := PointsByDist{Points: candidatePoints, ReferencePoint: startingPoint}

	sort.Sort(pd)

	for _, p := range pd.Points[:k] {
		fmt.Println(p)
	}
}
