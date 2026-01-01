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

func (b Bounds) Intersects(other Bounds) bool {
	/*
		Check if the two boxes touch or are separate, returns true if they touch, and false if they are separate,
		by doing an invariant
	*/
	return !(b.X > other.X+other.Width ||
		b.X+b.Width < other.X ||
		b.Y > other.Y+other.Height ||
		b.Y+b.Height < other.Y)
}

func (b Bounds) Contains(point Point) bool {
	return point.X >= b.X && point.X <= b.X+b.Width &&
		point.Y <= b.Y+b.Height && point.Y >= b.Y
}

func (n *Node) SubDivide() {
	x := n.Bounds.X
	y := n.Bounds.Y
	w := n.Bounds.X / 2
	h := n.Bounds.Y / 2
	//NW Child
	n.Children[0] = &Node{
		Bounds:   Bounds{X: x, Y: y, Width: w, Height: h},
		Capacity: n.Capacity,
	}
	//NE Child
	n.Children[1] = &Node{
		Bounds:   Bounds{X: x + w, Y: y, Width: w, Height: h},
		Capacity: n.Capacity,
	}
	//SW Child
	n.Children[2] = &Node{
		Bounds:   Bounds{X: x, Y: y - h, Width: w, Height: h},
		Capacity: n.Capacity,
	}
	//SE Child
	n.Children[3] = &Node{
		Bounds:   Bounds{X: x + w, Y: y - h, Width: w, Height: h},
		Capacity: n.Capacity,
	}
	for _, p := range n.Points {
		for i := 0; i < 4; i++ {
			if n.Children[i].InsertNode(p) {
				break
			}
		}
	}
	n.Points = nil

}

func (n *Node) InsertNode(point Point) bool {
	if n.Bounds.Contains(point) == false {
		return false
	}
	if len(n.Points) < n.Capacity && n.Children[0] == nil {
		n.Points = append(n.Points, point)
	}
	if n.Children[0] == nil {
		n.SubDivide()
	}
	for i := 0; i < 4; i++ {
		if n.Children[i].InsertNode(point) {
			return true
		}
	}
	return false
}

func (n *Node) Search(searchArea Bounds, resultPoints *[]Point) {
	if n.Bounds.Intersects(searchArea) == false {
		return
	}
	if n.Points != nil {
		for _, p := range n.Points {
			if searchArea.Contains(p) {
				*resultPoints = append(*resultPoints, p)
			}
		}

	} else {
		for i := 0; i < 4; i++ {
			n.Children[i].Search(searchArea, resultPoints)
		}
	}
}
