// Package spatial provides hierarchical spatial indexing using QuadTree data structure
// for efficient geometric queries and nearest-neighbor searches on 2D coordinates.
// It supports thread-safe insertion, deletion, searching, and k-nearest neighbor operations.
package spatial

import (
	"math"
	"sync"
)

// Point represents a 2D coordinate with optional associated data.
type Point struct {
	X, Y float64     // Coordinates in 2D space
	Data interface{} // Arbitrary data (e.g., vehicle ID, delivery ID)
}

// Bounds defines a rectangular region in 2D space.
type Bounds struct {
	X      float64 // Top-Left X coordinate
	Y      float64 // Top-Left Y coordinate
	Width  float64 // Width of the rectangular region
	Height float64 // Height of the rectangular region
}

// Node represents a single node in the QuadTree structure.
type Node struct {
	Bounds   Bounds   // The bounding box (X, Y, width, height) of the region covered by this node
	Points   []Point  // Points stored in this node (only populated in leaf nodes)
	Capacity int      // Maximum points before subdivision into child nodes
	Children [4]*Node // Four child nodes: [0]=NW, [1]=NE, [2]=SW, [3]=SE (nil until subdivided)
}

// QuadTree is a thread-safe hierarchical spatial index for 2D points.
// It recursively partitions a 2D space into quadrants for efficient spatial queries.
type QuadTree struct {
	Root *Node        // Root node of the tree
	Lock sync.RWMutex // Synchronizes concurrent access to the tree
}

// PointWithDistance is a helper struct for sorting points by distance.
type PointWithDistance struct {
	Point    Point   // The point in space
	Distance float64 // The computed distance from a reference point
}

// Intersects checks if this bounds region intersects with another bounds region.
// Returns true if the regions touch or overlap, false if they are completely separate.
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

// Contains checks if the given point is within this bounds region.
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

// Update atomically updates a point location from oldPoint to newPoint.
// Returns true if the update succeeded, false if the new point is out of bounds.
// If insertion fails, the old point is restored.
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

// Remove atomically removes a point from the QuadTree.
// Returns true if the point existed and was removed, false otherwise.
// This operation acquires a write lock for thread-safety.
func (qt *QuadTree) Remove(point Point) bool {
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	return qt.Root.RemoveNode(point)
}

// Insert atomically inserts a point into the QuadTree.
// Returns true if insertion succeeded, false if the point is out of bounds.
// This operation acquires a write lock for thread-safety.
func (qt *QuadTree) Insert(point Point) bool {
	qt.Lock.Lock()
	defer qt.Lock.Unlock()
	res := qt.Root.InsertNode(point)
	return res
}

// Search atomically retrieves all points within the given rectangular area.
// Returns a slice of points that intersect with the search bounds.
// This operation acquires a read lock allowing concurrent searches.
func (qt *QuadTree) Search(area Bounds) []Point {
	qt.Lock.RLock()
	defer qt.Lock.RUnlock()
	results := make([]Point, 0)
	qt.Root.SearchTree(area, &results)
	return results
}

// Distance calculates the Euclidean distance between two points.
func Distance(p1, p2 Point) float64 {
	dx := p2.X - p1.X
	dy := p2.Y - p1.Y
	return math.Sqrt(dx*dx + dy*dy)
}

// sortByDistance sorts points in ascending order by their distance field using insertion sort.
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

// KNearest finds the k nearest points to a target point in the QuadTree.
// Uses expanding search radius to locate candidates, then sorts by distance.
// Returns up to k points sorted by distance from the target.
// Returns empty slice if k <= 0 or tree is empty.
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
