package spatial

import "sync"

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

type QuadTree struct {
	Root *Node
	Lock sync.RWMutex
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
	w := n.Bounds.Width / 2
	h := n.Bounds.Height / 2
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
		Bounds:   Bounds{X: x, Y: y + h, Width: w, Height: h},
		Capacity: n.Capacity,
	}
	//SE Child
	n.Children[3] = &Node{
		Bounds:   Bounds{X: x + w, Y: y + h, Width: w, Height: h},
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

//Internal Function
func (n *Node) InsertNode(point Point) bool {
	if n.Bounds.Contains(point) == false {
		return false
	}
	if n.Children[0] != nil {
		for i := 0; i < 4; i++ {
			if n.Children[i].InsertNode(point) {
				return true
			}
		}
		return false
	}
	if len(n.Points) < n.Capacity && n.Children[0] == nil {
		n.Points = append(n.Points, point)
		return true
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

//Internal Function
func (n *Node) SearchTree(searchArea Bounds, resultPoints *[]Point) {

	if n == nil || !n.Bounds.Intersects(searchArea) {
		return
	}
	if n.Children[0] != nil {
		for i := 0; i < 4; i++ {
			n.Children[i].SearchTree(searchArea, resultPoints)
		}
		return
	}
	for _, p := range n.Points {
		if searchArea.Contains(p) {
			*resultPoints = append(*resultPoints, p)
		}
	}
}

func (qt *QuadTree) Insert(point Point) bool {
	/*
		Public Accessible API for Inserting New Points into the QuadTree,
		(Much Less Contention in Comparison to the Read Operations)
	*/
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	res := qt.Root.InsertNode(point)
	return res

}

func (qt *QuadTree) Search(area Bounds) []Point {
	/*
		Public Accessible API to search within the QuadTree
	*/
	qt.Lock.RLock()
	defer qt.Lock.RUnlock()
	results := make([]Point, 0)
	qt.Root.SearchTree(area, &results)
	return results
}
