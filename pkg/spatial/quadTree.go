package spatial

import (
	"math"
	"sync"
)

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

// PointWithDistance is a helper struct for sorting points by distance
type PointWithDistance struct {
	Point    Point
	Distance float64
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

// Internal Function for Inserting a Node
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

// Internal Function for Searching within the Tree
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

func (n *Node) RemoveNode(point Point) bool {
	if !n.Bounds.Contains(point) {
		return false
	}
	if n.Children[0] != nil { //If Node isnt a leaf node
		for i := 0; i < 4; i++ {
			if n.Children[i].RemoveNode(point) {
				return true
			}
		}
		return false
	}
	for i, exist := range n.Points {
		if exist.X == point.X && exist.Y == point.Y { //Switching the found value to the last, and slicing it, as order doesnt matter
			n.Points[i] = n.Points[len(n.Points)-1]
			n.Points = n.Points[:len(n.Points)-1]
			return true
		}
	}
	return false

}

func (qt *QuadTree) Update(oldPoint, newPoint Point) bool {
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	// Validate new point is within bounds before removing old point
	if !qt.Root.Bounds.Contains(newPoint) {
		return false
	}
	if qt.Root.RemoveNode(oldPoint) {
		if qt.Root.InsertNode(newPoint) {
			return true
		}
		//re-insert old point if new insert failed
		qt.Root.InsertNode(oldPoint)
		return false
	}
	return false
}

func (qt *QuadTree) Remove(point Point) bool {
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	return qt.Root.RemoveNode(point)
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

func Distance(p1, p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

func sortByDistance(points []PointWithDistance) {
	for i := 1; i < len(points); i++ {
		key := points[i]
		j := i - 1

		for j >= 0 && points[j].Distance > key.Distance {
			points[j+1] = points[j]
			j--
		}
		points[j+1] = key
	}
}

func (qt *QuadTree) KNearest(target Point, k int) []Point {
	if k <= 0 {
		return make([]Point, 0)
	}

	qt.Lock.RLock()
	defer qt.Lock.RUnlock()

	if qt.Root == nil {
		return make([]Point, 0)
	}

	maxPoints := k * 10
	if maxPoints > 100000 {
		maxPoints = 100000
	}

	// Start with initial search radius
	initialRadius := 10.0
	searchRadius := initialRadius
	maxRadius := math.Max(qt.Root.Bounds.Width, qt.Root.Bounds.Height) * 2
	var results []Point

	for searchRadius <= maxRadius {

		searchBounds := Bounds{
			X:      target.X - searchRadius,
			Y:      target.Y - searchRadius,
			Width:  searchRadius * 2,
			Height: searchRadius * 2,
		}

		results = make([]Point, 0)
		qt.Root.SearchTree(searchBounds, &results)

		if len(results) >= k {
			break
		}

		if searchRadius >= maxRadius {
			break
		}
		if len(results) > maxPoints {
			break
		}

		searchRadius *= 2
	}

	if len(results) == 0 {
		return make([]Point, 0)
	}

	pointsWithDist := make([]PointWithDistance, len(results))
	for i, p := range results {
		pointsWithDist[i] = PointWithDistance{
			Point:    p,
			Distance: Distance(target, p),
		}
	}

	sortByDistance(pointsWithDist)

	if len(pointsWithDist) > k {
		pointsWithDist = pointsWithDist[:k]
	}

	finalResults := make([]Point, len(pointsWithDist))
	for i, pd := range pointsWithDist {
		finalResults[i] = pd.Point
	}

	return finalResults
}
