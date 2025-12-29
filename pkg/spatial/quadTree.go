package main

type Point struct {
	X, Y float64
	Data interface{}
}
type Bounds struct {
	X      float64 // Top-Left X coordinate
	Y      float64 // Top-Left Y coordinate
	Width  float64
	Height float64
}

type Node struct {
	Bounds   Bounds //The X,Y height, width of the box
	Points   []Point
	Capacity int
	Children [4]*Node
}

func (b Bounds) Contains(point Point) bool {
	return point.X >= b.X && point.X <= b.X+b.Width &&
		point.Y <= b.Y+b.Height && point.Y >= b.Y
}
