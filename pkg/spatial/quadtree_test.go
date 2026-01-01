package spatial

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

// TestBoundsContains tests if a point is correctly identified as within bounds
func TestBoundsContains(t *testing.T) {
	tests := []struct {
		name     string
		bounds   Bounds
		point    Point
		expected bool
	}{
		{
			name:     "point inside bounds",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 5, Y: 5, Data: nil},
			expected: true,
		},
		{
			name:     "point on left edge",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 0, Y: 5, Data: nil},
			expected: true,
		},
		{
			name:     "point on right edge",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 10, Y: 5, Data: nil},
			expected: true,
		},
		{
			name:     "point on top edge",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 5, Y: 0, Data: nil},
			expected: true,
		},
		{
			name:     "point on bottom edge",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 5, Y: 10, Data: nil},
			expected: true,
		},
		{
			name:     "point outside bounds left",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: -1, Y: 5, Data: nil},
			expected: false,
		},
		{
			name:     "point outside bounds right",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 11, Y: 5, Data: nil},
			expected: false,
		},
		{
			name:     "point outside bounds top",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 5, Y: -1, Data: nil},
			expected: false,
		},
		{
			name:     "point outside bounds bottom",
			bounds:   Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			point:    Point{X: 5, Y: 11, Data: nil},
			expected: false,
		},
		{
			name:     "negative coordinates inside bounds",
			bounds:   Bounds{X: -10, Y: -10, Width: 20, Height: 20},
			point:    Point{X: 0, Y: 0, Data: nil},
			expected: true,
		},
		{
			name:     "negative coordinates outside bounds",
			bounds:   Bounds{X: -10, Y: -10, Width: 10, Height: 10},
			point:    Point{X: 5, Y: 5, Data: nil},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bounds.Contains(tt.point)
			if result != tt.expected {
				t.Errorf("Contains() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestBoundsIntersects tests boundary intersection detection
func TestBoundsIntersects(t *testing.T) {
	tests := []struct {
		name     string
		bounds1  Bounds
		bounds2  Bounds
		expected bool
	}{
		{
			name:     "overlapping bounds",
			bounds1:  Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			bounds2:  Bounds{X: 5, Y: 5, Width: 10, Height: 10},
			expected: true,
		},
		{
			name:     "touching bounds on edge",
			bounds1:  Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			bounds2:  Bounds{X: 10, Y: 0, Width: 10, Height: 10},
			expected: true,
		},
		{
			name:     "separated bounds",
			bounds1:  Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			bounds2:  Bounds{X: 20, Y: 20, Width: 10, Height: 10},
			expected: false,
		},
		{
			name:     "contained bounds",
			bounds1:  Bounds{X: 0, Y: 0, Width: 20, Height: 20},
			bounds2:  Bounds{X: 5, Y: 5, Width: 5, Height: 5},
			expected: true,
		},
		{
			name:     "adjacent vertically",
			bounds1:  Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			bounds2:  Bounds{X: 0, Y: 10, Width: 10, Height: 10},
			expected: true,
		},
		{
			name:     "separated horizontally",
			bounds1:  Bounds{X: 0, Y: 0, Width: 10, Height: 10},
			bounds2:  Bounds{X: 11, Y: 0, Width: 10, Height: 10},
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.bounds1.Intersects(tt.bounds2)
			if result != tt.expected {
				t.Errorf("Intersects() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestNodeSubDivide tests that a node correctly subdivides into 4 children
func TestNodeSubDivide(t *testing.T) {
	node := &Node{
		Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
		Capacity: 4,
	}

	node.SubDivide()

	if node.Children[0] == nil || node.Children[1] == nil || node.Children[2] == nil || node.Children[3] == nil {
		t.Error("SubDivide() failed to create all 4 children")
	}

	// Check NW child
	if node.Children[0].Bounds.X != 0 || node.Children[0].Bounds.Y != 0 {
		t.Error("NW child bounds incorrect")
	}

	// Check NE child
	if node.Children[1].Bounds.X != 50 || node.Children[1].Bounds.Y != 0 {
		t.Error("NE child bounds incorrect")
	}

	// Check SW child
	if node.Children[2].Bounds.X != 0 || node.Children[2].Bounds.Y != 50 {
		t.Error("SW child bounds incorrect")
	}

	// Check SE child
	if node.Children[3].Bounds.X != 50 || node.Children[3].Bounds.Y != 50 {
		t.Error("SE child bounds incorrect")
	}

	// Verify all children have correct dimensions
	for i := 0; i < 4; i++ {
		if node.Children[i].Bounds.Width != 50 || node.Children[i].Bounds.Height != 50 {
			t.Errorf("Child %d has incorrect dimensions", i)
		}
		if node.Children[i].Capacity != node.Capacity {
			t.Errorf("Child %d has incorrect capacity", i)
		}
	}
}

// TestNodeSubDivideWithPoints tests that points are correctly distributed during subdivision
func TestNodeSubDivideWithPoints(t *testing.T) {
	node := &Node{
		Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
		Capacity: 2,
		Points: []Point{
			{X: 25, Y: 25, Data: "NW"},
			{X: 75, Y: 25, Data: "NE"},
			{X: 25, Y: 75, Data: "SW"},
			{X: 75, Y: 75, Data: "SE"},
		},
	}

	node.SubDivide()

	// After subdivision, original points should be cleared
	if len(node.Points) != 0 {
		t.Errorf("Parent node points should be empty after subdivision, got %d", len(node.Points))
	}

	// Check that points were distributed to children
	totalPointsInChildren := 0
	for i := 0; i < 4; i++ {
		totalPointsInChildren += len(node.Children[i].Points)
	}

	if totalPointsInChildren != 4 {
		t.Errorf("Expected 4 points in children, got %d", totalPointsInChildren)
	}
}

// TestQuadTreeInitialization tests QuadTree creation
func TestQuadTreeInitialization(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	if qt.Root == nil {
		t.Error("QuadTree Root is nil")
	}

	if qt.Root.Bounds.X != 0 || qt.Root.Bounds.Y != 0 {
		t.Error("QuadTree Root bounds incorrect")
	}
}

// TestQuadTreeInsertSingle tests inserting a single point
func TestQuadTreeInsertSingle(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	point := Point{X: 50, Y: 50, Data: "test"}
	result := qt.Insert(point)

	if !result {
		t.Error("Insert() failed to insert point")
	}

	if len(qt.Root.Points) != 1 {
		t.Errorf("Expected 1 point in root, got %d", len(qt.Root.Points))
	}

	if qt.Root.Points[0].X != 50 || qt.Root.Points[0].Y != 50 {
		t.Error("Inserted point has incorrect coordinates")
	}
}

// TestQuadTreeInsertMultiple tests inserting multiple points
func TestQuadTreeInsertMultiple(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
		{X: 40, Y: 40, Data: "p4"},
	}

	for _, p := range points {
		result := qt.Insert(p)
		if !result {
			t.Errorf("Failed to insert point %v", p)
		}
	}

	if len(qt.Root.Points) != 4 {
		t.Errorf("Expected 4 points in root, got %d", len(qt.Root.Points))
	}
}

// TestQuadTreeInsertExceedsCapacity tests insertion after capacity is exceeded
func TestQuadTreeInsertExceedsCapacity(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 2,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
		{X: 40, Y: 40, Data: "p4"},
	}

	for i, p := range points {
		result := qt.Insert(p)
		if !result {
			t.Errorf("Failed to insert point %d", i)
		}
	}

	// After exceeding capacity, root should be subdivided
	if qt.Root.Children[0] == nil {
		t.Error("Root should be subdivided after exceeding capacity")
	}

	// Root should have no points after subdivision
	if len(qt.Root.Points) != 0 {
		t.Error("Root points should be cleared after subdivision")
	}
}

// TestQuadTreeInsertOutOfBounds tests inserting points outside bounds
func TestQuadTreeInsertOutOfBounds(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	point := Point{X: 150, Y: 150, Data: "out of bounds"}
	result := qt.Insert(point)

	if result {
		t.Error("Insert() should fail for out of bounds point")
	}

	if len(qt.Root.Points) != 0 {
		t.Error("Out of bounds point should not be inserted")
	}
}

// TestQuadTreeSearchBasic tests basic search functionality
func TestQuadTreeSearchBasic(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
		{X: 40, Y: 40, Data: "p4"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Search for points in specific area
	searchArea := Bounds{X: 0, Y: 0, Width: 25, Height: 25}
	results := qt.Search(searchArea)

	if len(results) != 2 {
		t.Errorf("Expected 2 results, got %d", len(results))
	}
}

// TestQuadTreeSearchNoResults tests search that yields no results
func TestQuadTreeSearchNoResults(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Search in area with no points
	searchArea := Bounds{X: 50, Y: 50, Width: 20, Height: 20}
	results := qt.Search(searchArea)

	if len(results) != 0 {
		t.Errorf("Expected 0 results, got %d", len(results))
	}
}

// TestQuadTreeSearchAll tests searching entire tree
func TestQuadTreeSearchAll(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	numPoints := 20
	for i := 0; i < numPoints; i++ {
		x := float64(i*5) + 5
		y := float64(i*5) + 5
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", i)})
	}

	// Search entire bounds
	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != numPoints {
		t.Errorf("Expected %d results, got %d", numPoints, len(results))
	}
}

// TestQuadTreeSearchAfterSubdivision tests search functionality after tree subdivision
func TestQuadTreeSearchAfterSubdivision(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 2,
		},
	}

	// Insert points to trigger subdivision
	points := []Point{
		{X: 10, Y: 10, Data: "NW"},
		{X: 80, Y: 10, Data: "NE"},
		{X: 10, Y: 80, Data: "SW"},
		{X: 80, Y: 80, Data: "SE"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Verify subdivision occurred
	if qt.Root.Children[0] == nil {
		t.Fatal("Tree should be subdivided")
	}

	// Search in different regions
	nwArea := Bounds{X: 0, Y: 0, Width: 50, Height: 50}
	nwResults := qt.Search(nwArea)
	if len(nwResults) != 1 {
		t.Errorf("NW search: expected 1 result, got %d", len(nwResults))
	}

	neArea := Bounds{X: 50, Y: 0, Width: 50, Height: 50}
	neResults := qt.Search(neArea)
	if len(neResults) != 1 {
		t.Errorf("NE search: expected 1 result, got %d", len(neResults))
	}
}

// TestQuadTreeConcurrentInsert tests concurrent insert operations
func TestQuadTreeConcurrentInsert(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 1000, Height: 1000},
			Capacity: 10,
		},
	}

	var wg sync.WaitGroup
	numGoroutines := 10
	pointsPerGoroutine := 50

	wg.Add(numGoroutines)
	for g := 0; g < numGoroutines; g++ {
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < pointsPerGoroutine; i++ {
				x := float64(goroutineID*100 + i*2)
				y := float64(goroutineID*100 + i*2)
				point := Point{X: x, Y: y, Data: fmt.Sprintf("g%d_p%d", goroutineID, i)}
				qt.Insert(point)
			}
		}(g)
	}

	wg.Wait()

	// Search to verify all points were inserted
	searchArea := Bounds{X: 0, Y: 0, Width: 1000, Height: 1000}
	results := qt.Search(searchArea)

	expectedPoints := numGoroutines * pointsPerGoroutine
	if len(results) != expectedPoints {
		t.Errorf("Expected %d points, got %d", expectedPoints, len(results))
	}
}

// TestQuadTreeConcurrentSearchInsert tests concurrent search and insert operations
func TestQuadTreeConcurrentSearchInsert(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 1000, Height: 1000},
			Capacity: 10,
		},
	}

	// Pre-populate with some points
	for i := 0; i < 100; i++ {
		x := float64(i * 10)
		y := float64(i * 10)
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("initial_%d", i)})
	}

	var wg sync.WaitGroup
	numReaders := 5
	numWriters := 5

	// Reader goroutines
	for r := 0; r < numReaders; r++ {
		wg.Add(1)
		go func(readerID int) {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				searchArea := Bounds{X: float64(readerID * 100), Y: float64(readerID * 100), Width: 200, Height: 200}
				_ = qt.Search(searchArea)
			}
		}(r)
	}

	// Writer goroutines
	for w := 0; w < numWriters; w++ {
		wg.Add(1)
		go func(writerID int) {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				x := float64(writerID*200+i) + 500
				y := float64(writerID*200+i) + 500
				qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("w%d_p%d", writerID, i)})
			}
		}(w)
	}

	wg.Wait()
}

// TestQuadTreeSearchWithData tests that point data is preserved
func TestQuadTreeSearchWithData(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	testData := []interface{}{"location1", 42, 3.14, true}
	for i, data := range testData {
		x := float64(i*20) + 10
		y := float64(i*20) + 10
		qt.Insert(Point{X: x, Y: y, Data: data})
	}

	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != len(testData) {
		t.Errorf("Expected %d results, got %d", len(testData), len(results))
	}

	for _, result := range results {
		if result.Data == nil {
			t.Error("Result data should not be nil")
		}
	}
}

// TestQuadTreeSearchPrecision tests search with floating point precision
func TestQuadTreeSearchPrecision(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	// Insert points with fractional coordinates
	points := []Point{
		{X: 10.5, Y: 10.5, Data: "p1"},
		{X: 20.75, Y: 20.25, Data: "p2"},
		{X: 99.99, Y: 99.99, Data: "p3"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	searchArea := Bounds{X: 10.0, Y: 10.0, Width: 20.0, Height: 20.0}
	results := qt.Search(searchArea)

	if len(results) != 2 {
		t.Errorf("Expected 2 results with fractional coordinates, got %d", len(results))
	}
}

// TestQuadTreeDeepNesting tests tree with deep nesting/subdivision
func TestQuadTreeDeepNesting(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 1024, Height: 1024},
			Capacity: 1,
		},
	}

	// Insert enough points to force deep subdivision
	for i := 0; i < 100; i++ {
		x := float64(i%10) * 10
		y := float64(i/10) * 10
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", i)})
	}

	// Verify all points can be retrieved
	searchArea := Bounds{X: 0, Y: 0, Width: 1024, Height: 1024}
	results := qt.Search(searchArea)

	if len(results) != 100 {
		t.Errorf("Expected 100 results in deep tree, got %d", len(results))
	}
}

// TestQuadTreeNegativeCoordinates tests tree with negative coordinates
func TestQuadTreeNegativeCoordinates(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: -100, Y: -100, Width: 200, Height: 200},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: -50, Y: -50, Data: "neg1"},
		{X: 0, Y: 0, Data: "origin"},
		{X: 50, Y: 50, Data: "pos1"},
		{X: -25, Y: 75, Data: "mixed"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	searchArea := Bounds{X: -100, Y: -100, Width: 200, Height: 200}
	results := qt.Search(searchArea)

	if len(results) != 4 {
		t.Errorf("Expected 4 results with negative coordinates, got %d", len(results))
	}
}

// BenchmarkInsert benchmarks insertion performance
func BenchmarkInsert(b *testing.B) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 10,
		},
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		x := float64(i%100) * 100
		y := float64((i/100)%100) * 100
		qt.Insert(Point{X: x, Y: y, Data: nil})
	}
}

// TestQuadTreeRemoveBasic tests removing a single point
func TestQuadTreeRemoveBasic(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	point := Point{X: 50, Y: 50, Data: "test"}
	qt.Insert(point)

	if len(qt.Root.Points) != 1 {
		t.Errorf("Expected 1 point after insert, got %d", len(qt.Root.Points))
	}

	// Remove the point
	result := qt.Remove(point)

	if !result {
		t.Error("Remove() failed to remove existing point")
	}

	if len(qt.Root.Points) != 0 {
		t.Errorf("Expected 0 points after removal, got %d", len(qt.Root.Points))
	}
}

// TestQuadTreeRemoveNonexistent tests removing a point that doesn't exist
func TestQuadTreeRemoveNonexistent(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	point := Point{X: 50, Y: 50, Data: "test"}
	result := qt.Remove(point)

	if result {
		t.Error("Remove() should fail for non-existent point")
	}
}

// TestQuadTreeRemoveMultiple tests removing multiple points from a leaf node
func TestQuadTreeRemoveMultiple(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Remove first point
	if !qt.Remove(points[0]) {
		t.Error("Failed to remove first point")
	}

	if len(qt.Root.Points) != 2 {
		t.Errorf("Expected 2 points after removal, got %d", len(qt.Root.Points))
	}

	// Remove second point
	if !qt.Remove(points[1]) {
		t.Error("Failed to remove second point")
	}

	if len(qt.Root.Points) != 1 {
		t.Errorf("Expected 1 point after removal, got %d", len(qt.Root.Points))
	}
}

// TestQuadTreeRemoveAfterSubdivision tests removing points after tree subdivision
func TestQuadTreeRemoveAfterSubdivision(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 2,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "NW"},
		{X: 80, Y: 10, Data: "NE"},
		{X: 10, Y: 80, Data: "SW"},
		{X: 80, Y: 80, Data: "SE"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Verify tree was subdivided
	if qt.Root.Children[0] == nil {
		t.Fatal("Tree should be subdivided")
	}

	// Remove a point from a child node
	if !qt.Remove(points[0]) {
		t.Error("Failed to remove point from subdivided tree")
	}

	// Verify point was removed
	searchArea := Bounds{X: 0, Y: 0, Width: 50, Height: 50}
	results := qt.Search(searchArea)

	if len(results) != 0 {
		t.Errorf("Expected 0 results in NW quadrant, got %d", len(results))
	}
}

// TestQuadTreeRemoveOutOfBounds tests removing an out-of-bounds point
func TestQuadTreeRemoveOutOfBounds(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	qt.Insert(Point{X: 50, Y: 50, Data: "test"})

	// Try to remove a point outside bounds
	outOfBounds := Point{X: 150, Y: 150, Data: "test"}
	result := qt.Remove(outOfBounds)

	if result {
		t.Error("Remove() should fail for out-of-bounds point")
	}
}

// TestQuadTreeUpdateBasic tests updating a point's coordinates
func TestQuadTreeUpdateBasic(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	oldPoint := Point{X: 50, Y: 50, Data: "test"}
	qt.Insert(oldPoint)

	if len(qt.Root.Points) != 1 {
		t.Fatal("Insert failed")
	}

	// Update the point to new coordinates
	newPoint := Point{X: 75, Y: 75, Data: "updated"}
	result := qt.Update(oldPoint, newPoint)

	if !result {
		t.Error("Update() failed")
	}

	// Verify old point is gone
	searchOld := Bounds{X: 45, Y: 45, Width: 10, Height: 10}
	resultsOld := qt.Search(searchOld)

	if len(resultsOld) != 0 {
		t.Error("Old point should not exist after update")
	}

	// Verify new point exists
	searchNew := Bounds{X: 70, Y: 70, Width: 10, Height: 10}
	resultsNew := qt.Search(searchNew)

	if len(resultsNew) != 1 {
		t.Errorf("Expected 1 result at new coordinates, got %d", len(resultsNew))
	}

	if resultsNew[0].Data != "updated" {
		t.Error("Updated point should have new data")
	}
}

// TestQuadTreeUpdateNonexistentPoint tests updating a point that doesn't exist
func TestQuadTreeUpdateNonexistentPoint(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	oldPoint := Point{X: 50, Y: 50, Data: "test"}
	newPoint := Point{X: 75, Y: 75, Data: "updated"}

	result := qt.Update(oldPoint, newPoint)

	if result {
		t.Error("Update() should fail for non-existent point")
	}

	// Verify new point was not inserted
	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != 0 {
		t.Errorf("Expected no results, got %d", len(results))
	}
}

// TestQuadTreeUpdateToOutOfBounds tests updating a point outside the bounds
func TestQuadTreeUpdateToOutOfBounds(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	oldPoint := Point{X: 50, Y: 50, Data: "test"}
	qt.Insert(oldPoint)

	// Try to update to out-of-bounds coordinates
	outOfBoundsPoint := Point{X: 150, Y: 150, Data: "out"}
	result := qt.Update(oldPoint, outOfBoundsPoint)

	if result {
		t.Error("Update() should fail when moving point out of bounds")
	}

	// Old point should still exist
	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != 1 {
		t.Error("Original point should still exist after failed update")
	}
}

// TestQuadTreeUpdateMultiplePoints tests updating multiple points sequentially
func TestQuadTreeUpdateMultiplePoints(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Update multiple points
	newPoints := []Point{
		{X: 15, Y: 15, Data: "p1_updated"},
		{X: 25, Y: 25, Data: "p2_updated"},
		{X: 35, Y: 35, Data: "p3_updated"},
	}

	for i, newP := range newPoints {
		if !qt.Update(points[i], newP) {
			t.Errorf("Failed to update point %d", i)
		}
	}

	// Verify all updates
	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != 3 {
		t.Errorf("Expected 3 results after updates, got %d", len(results))
	}
}

// TestQuadTreeUpdateAfterSubdivision tests updating points after tree subdivision
func TestQuadTreeUpdateAfterSubdivision(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 2,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "NW"},
		{X: 80, Y: 10, Data: "NE"},
		{X: 10, Y: 80, Data: "SW"},
		{X: 80, Y: 80, Data: "SE"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Verify tree was subdivided
	if qt.Root.Children[0] == nil {
		t.Fatal("Tree should be subdivided")
	}

	// Update a point to different quadrant
	oldPoint := points[0] // In NW
	newPoint := Point{X: 80, Y: 80, Data: "moved_to_SE"}

	if !qt.Update(oldPoint, newPoint) {
		t.Error("Update() should succeed")
	}

	// Verify old quadrant has one less point
	nwArea := Bounds{X: 0, Y: 0, Width: 50, Height: 50}
	nwResults := qt.Search(nwArea)

	if len(nwResults) != 0 {
		t.Errorf("NW quadrant should have 0 points, got %d", len(nwResults))
	}

	// Verify new quadrant has the updated point
	seArea := Bounds{X: 50, Y: 50, Width: 50, Height: 50}
	seResults := qt.Search(seArea)

	if len(seResults) != 2 {
		t.Errorf("SE quadrant should have 2 points, got %d", len(seResults))
	}
}

// TestQuadTreeRemoveFromDifferentQuadrants tests removing points from various quadrants
func TestQuadTreeRemoveFromDifferentQuadrants(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 2,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "NW"},
		{X: 80, Y: 10, Data: "NE"},
		{X: 10, Y: 80, Data: "SW"},
		{X: 80, Y: 80, Data: "SE"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Remove from each quadrant
	for i, p := range points {
		if !qt.Remove(p) {
			t.Errorf("Failed to remove point %d", i)
		}
	}

	// Verify all points are removed
	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != 0 {
		t.Errorf("Expected all points removed, got %d", len(results))
	}
}

// TestQuadTreeRemoveSameName tests removing points with same data value
func TestQuadTreeRemoveSameName(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	// Insert points with same data but different coordinates
	points := []Point{
		{X: 10, Y: 10, Data: "duplicate"},
		{X: 20, Y: 20, Data: "duplicate"},
		{X: 30, Y: 30, Data: "unique"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Remove should match by coordinates, not data
	if !qt.Remove(points[0]) {
		t.Error("Failed to remove first point")
	}

	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != 2 {
		t.Errorf("Expected 2 results after removal, got %d", len(results))
	}
}

// TestQuadTreeUpdatePreservesDataInSearchArea tests update preserves point data
func TestQuadTreeUpdatePreservesDataInSearchArea(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	oldPoint := Point{X: 10, Y: 10, Data: "original_data"}
	qt.Insert(oldPoint)

	newPoint := Point{X: 30, Y: 30, Data: "new_data"}
	qt.Update(oldPoint, newPoint)

	// Search at old location
	oldSearchArea := Bounds{X: 5, Y: 5, Width: 10, Height: 10}
	oldResults := qt.Search(oldSearchArea)

	if len(oldResults) != 0 {
		t.Error("Old location should have no points")
	}

	// Search at new location
	newSearchArea := Bounds{X: 25, Y: 25, Width: 10, Height: 10}
	newResults := qt.Search(newSearchArea)

	if len(newResults) != 1 {
		t.Error("New location should have the point")
	}

	if newResults[0].Data != "new_data" {
		t.Error("Point should have updated data")
	}
}

// TestQuadTreeConcurrentUpdateRemove tests concurrent update and remove operations
func TestQuadTreeConcurrentUpdateRemove(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 1000, Height: 1000},
			Capacity: 10,
		},
	}

	// Pre-populate
	for i := 0; i < 100; i++ {
		x := float64(i%10) * 100
		y := float64((i/10)%10) * 100
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", i)})
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	// Half update, half remove
	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			if goroutineID%2 == 0 {
				// Update operations
				for i := 0; i < 10; i++ {
					oldX := float64((goroutineID*10 + i) % 100)
					oldY := float64((goroutineID*10 + i) / 100)
					oldP := Point{X: oldX, Y: oldY, Data: fmt.Sprintf("p%d", goroutineID*10+i)}

					newX := oldX + 500
					newY := oldY + 500
					newP := Point{X: newX, Y: newY, Data: fmt.Sprintf("updated_%d", goroutineID)}
					qt.Update(oldP, newP)
				}
			} else {
				// Remove operations
				for i := 0; i < 10; i++ {
					x := float64((goroutineID*10 + i) % 100)
					y := float64((goroutineID*10 + i) / 100)
					qt.Remove(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", goroutineID*10+i)})
				}
			}
		}(g)
	}

	wg.Wait()
}

// BenchmarkRemove benchmarks removal performance
func BenchmarkRemove(b *testing.B) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 10,
		},
	}

	// Pre-populate
	points := make([]Point, b.N)
	for i := 0; i < b.N; i++ {
		x := float64(i%100) * 100
		y := float64((i/100)%100) * 100
		p := Point{X: x, Y: y, Data: nil}
		points[i] = p
		qt.Insert(p)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		qt.Remove(points[i])
	}
}

// BenchmarkUpdate benchmarks update performance
func BenchmarkUpdate(b *testing.B) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 10,
		},
	}

	// Pre-populate
	points := make([]Point, b.N)
	for i := 0; i < b.N; i++ {
		x := float64(i%100) * 100
		y := float64((i/100)%100) * 100
		p := Point{X: x, Y: y, Data: nil}
		points[i] = p
		qt.Insert(p)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		newX := float64((i+5000)%100) * 100
		newY := float64(((i+5000)/100)%100) * 100
		qt.Update(points[i], Point{X: newX, Y: newY, Data: nil})
	}
}

// TestQuadTreeEmptySearch tests searching an empty tree
func TestQuadTreeEmptySearch(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	searchArea := Bounds{X: 0, Y: 0, Width: 100, Height: 100}
	results := qt.Search(searchArea)

	if len(results) != 0 {
		t.Errorf("Expected 0 results from empty tree, got %d", len(results))
	}

	if results == nil {
		t.Error("Search should return empty slice, not nil")
	}
}

// TestBoundsIntersectsSymmetry tests that intersection is symmetric
func TestBoundsIntersectsSymmetry(t *testing.T) {
	bounds1 := Bounds{X: 0, Y: 0, Width: 10, Height: 10}
	bounds2 := Bounds{X: 5, Y: 5, Width: 10, Height: 10}

	if bounds1.Intersects(bounds2) != bounds2.Intersects(bounds1) {
		t.Error("Intersection should be symmetric")
	}
}

// TestQuadTreeLargeScale tests insertion and search of large number of points
func TestQuadTreeLargeScale(t *testing.T) {
	if testing.Short() {
		t.Skip("Skipping large scale test in short mode")
	}

	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 20,
		},
	}

	// Insert 10000 points
	for i := 0; i < 10000; i++ {
		x := float64(i%100) * 100
		y := float64((i/100)%100) * 100
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", i)})
	}

	// Search various regions
	searchTests := []Bounds{
		{X: 0, Y: 0, Width: 1000, Height: 1000},
		{X: 2500, Y: 2500, Width: 2000, Height: 2000},
		{X: 9000, Y: 9000, Width: 1000, Height: 1000},
	}

	for _, searchArea := range searchTests {
		results := qt.Search(searchArea)
		if len(results) == 0 {
			t.Errorf("Search area %+v should have results", searchArea)
		}
	}
}

// TestBoundaryFloatOperations tests edge cases with floating point operations
func TestBoundaryFloatOperations(t *testing.T) {
	bounds := Bounds{X: 0, Y: 0, Width: math.MaxFloat64 / 2, Height: math.MaxFloat64 / 2}
	point := Point{X: math.MaxFloat64 / 4, Y: math.MaxFloat64 / 4, Data: nil}

	// Should not panic
	_ = bounds.Contains(point)
	_ = bounds.Intersects(Bounds{X: 0, Y: 0, Width: 100, Height: 100})
}

// TestDistanceCalculation tests the distance helper function
func TestDistanceCalculation(t *testing.T) {
	tests := []struct {
		name     string
		p1       Point
		p2       Point
		expected float64
	}{
		{
			name:     "same point",
			p1:       Point{X: 0, Y: 0, Data: nil},
			p2:       Point{X: 0, Y: 0, Data: nil},
			expected: 0,
		},
		{
			name:     "3-4-5 triangle",
			p1:       Point{X: 0, Y: 0, Data: nil},
			p2:       Point{X: 3, Y: 4, Data: nil},
			expected: 5,
		},
		{
			name:     "unit distance",
			p1:       Point{X: 0, Y: 0, Data: nil},
			p2:       Point{X: 1, Y: 0, Data: nil},
			expected: 1,
		},
		{
			name:     "negative coordinates",
			p1:       Point{X: -1, Y: -1, Data: nil},
			p2:       Point{X: 1, Y: 1, Data: nil},
			expected: math.Sqrt(8), // sqrt((1-(-1))^2 + (1-(-1))^2) = sqrt(8)
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := Distance(tt.p1, tt.p2)
			if math.Abs(result-tt.expected) > 1e-9 {
				t.Errorf("Distance() = %v, want %v", result, tt.expected)
			}
		})
	}
}

// TestDistanceSymmetry tests that Distance is symmetric
func TestDistanceSymmetry(t *testing.T) {
	p1 := Point{X: 10, Y: 20, Data: nil}
	p2 := Point{X: 30, Y: 40, Data: nil}

	d1 := Distance(p1, p2)
	d2 := Distance(p2, p1)

	if math.Abs(d1-d2) > 1e-9 {
		t.Errorf("Distance should be symmetric: d(p1,p2)=%v, d(p2,p1)=%v", d1, d2)
	}
}

// TestKNearestBasic tests basic k-nearest neighbor search
func TestKNearestBasic(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	// Insert points in a grid
	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
		{X: 40, Y: 40, Data: "p4"},
		{X: 50, Y: 50, Data: "p5"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Find 3 nearest to (15, 15)
	target := Point{X: 15, Y: 15, Data: nil}
	result := qt.KNearest(target, 3)

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	// Verify they're sorted by distance
	for i := 0; i < len(result)-1; i++ {
		d1 := Distance(target, result[i])
		d2 := Distance(target, result[i+1])
		if d1 > d2 {
			t.Errorf("Results not sorted by distance: d[%d]=%v > d[%d]=%v", i, d1, i+1, d2)
		}
	}
}

// TestKNearestZeroK tests with k=0
func TestKNearestZeroK(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	qt.Insert(Point{X: 50, Y: 50, Data: "test"})

	result := qt.KNearest(Point{X: 50, Y: 50, Data: nil}, 0)

	if len(result) != 0 {
		t.Errorf("Expected 0 results for k=0, got %d", len(result))
	}
}

// TestKNearestNegativeK tests with negative k
func TestKNearestNegativeK(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	qt.Insert(Point{X: 50, Y: 50, Data: "test"})

	result := qt.KNearest(Point{X: 50, Y: 50, Data: nil}, -5)

	if len(result) != 0 {
		t.Errorf("Expected 0 results for negative k, got %d", len(result))
	}
}

// TestKNearestEmptyTree tests k-nearest on empty tree
func TestKNearestEmptyTree(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	result := qt.KNearest(Point{X: 50, Y: 50, Data: nil}, 10)

	if len(result) != 0 {
		t.Errorf("Expected 0 results from empty tree, got %d", len(result))
	}
}

// TestKNearestMoreThanAvailable tests when k > available points
func TestKNearestMoreThanAvailable(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Request 10 nearest, but only 3 exist
	result := qt.KNearest(Point{X: 15, Y: 15, Data: nil}, 10)

	if len(result) != 3 {
		t.Errorf("Expected 3 results (all available), got %d", len(result))
	}
}

// TestKNearestExactDistance tests points at exact distances
func TestKNearestExactDistance(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 200, Height: 200},
			Capacity: 4,
		},
	}

	// Insert points at known distances from origin
	// Distance 5: (3, 4)
	// Distance 5: (4, 3)
	// Distance 10: (6, 8)
	qt.Insert(Point{X: 3, Y: 4, Data: "d5_a"})
	qt.Insert(Point{X: 4, Y: 3, Data: "d5_b"})
	qt.Insert(Point{X: 6, Y: 8, Data: "d10"})

	target := Point{X: 0, Y: 0, Data: nil}
	result := qt.KNearest(target, 2)

	if len(result) != 2 {
		t.Fatalf("Expected 2 results, got %d", len(result))
	}

	// Both should be at distance 5
	for i, p := range result {
		dist := Distance(target, p)
		if math.Abs(dist-5) > 1e-9 {
			t.Errorf("Result %d should be at distance 5, got %v", i, dist)
		}
	}
}

// TestKNearestTargetInsideTree tests k-nearest when target is in the tree
func TestKNearestTargetInsideTree(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	// Insert the target itself plus neighbors
	target := Point{X: 50, Y: 50, Data: "target"}
	qt.Insert(target)
	qt.Insert(Point{X: 51, Y: 50, Data: "neighbor1"})
	qt.Insert(Point{X: 50, Y: 51, Data: "neighbor2"})
	qt.Insert(Point{X: 10, Y: 10, Data: "far"})

	result := qt.KNearest(target, 2)

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}

	// First result should be target itself (distance 0)
	if Distance(target, result[0]) > 1e-9 {
		t.Error("Nearest point should be the target itself")
	}
}

// TestKNearestAfterSubdivision tests k-nearest after tree subdivision
func TestKNearestAfterSubdivision(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 2,
		},
	}

	// Insert enough points to trigger subdivision
	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
		{X: 40, Y: 40, Data: "p4"},
		{X: 50, Y: 50, Data: "p5"},
		{X: 60, Y: 60, Data: "p6"},
		{X: 70, Y: 70, Data: "p7"},
		{X: 80, Y: 80, Data: "p8"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	// Verify tree is subdivided
	if qt.Root.Children[0] == nil {
		t.Fatal("Tree should be subdivided")
	}

	// Find k nearest from middle of tree
	target := Point{X: 45, Y: 45, Data: nil}
	result := qt.KNearest(target, 3)

	if len(result) != 3 {
		t.Errorf("Expected 3 results, got %d", len(result))
	}

	// Verify ordering by distance
	for i := 0; i < len(result)-1; i++ {
		d1 := Distance(target, result[i])
		d2 := Distance(target, result[i+1])
		if d1 > d2 {
			t.Errorf("Results not sorted: d[%d]=%v > d[%d]=%v", i, d1, i+1, d2)
		}
	}
}

// TestKNearestNegativeCoordinates tests k-nearest with negative coordinates
func TestKNearestNegativeCoordinates(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: -100, Y: -100, Width: 200, Height: 200},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: -50, Y: -50, Data: "p1"},
		{X: 0, Y: 0, Data: "p2"},
		{X: 50, Y: 50, Data: "p3"},
		{X: -25, Y: 25, Data: "p4"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	target := Point{X: 0, Y: 0, Data: nil}
	result := qt.KNearest(target, 2)

	if len(result) != 2 {
		t.Errorf("Expected 2 results, got %d", len(result))
	}

	// Verify sorted by distance
	for i := 0; i < len(result)-1; i++ {
		d1 := Distance(target, result[i])
		d2 := Distance(target, result[i+1])
		if d1 > d2 {
			t.Errorf("Not sorted by distance")
		}
	}
}

// TestKNearestLargeDataset tests k-nearest on larger dataset
func TestKNearestLargeDataset(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 1000, Height: 1000},
			Capacity: 10,
		},
	}

	// Insert 100 random-ish points
	for i := 0; i < 100; i++ {
		x := float64(i%10) * 100
		y := float64((i/10)%10) * 100
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", i)})
	}

	target := Point{X: 450, Y: 450, Data: nil}
	k := 10
	result := qt.KNearest(target, k)

	if len(result) != k {
		t.Errorf("Expected %d results, got %d", k, len(result))
	}

	// Verify all results are sorted by distance
	for i := 0; i < len(result)-1; i++ {
		d1 := Distance(target, result[i])
		d2 := Distance(target, result[i+1])
		if d1 > d2+1e-9 { // Allow small floating point error
			t.Errorf("Not sorted at index %d: d[%d]=%v > d[%d]=%v", i, i, d1, i+1, d2)
		}
	}
}

// TestKNearestConsistency tests that KNearest returns consistent results
func TestKNearestConsistency(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 100, Height: 100},
			Capacity: 4,
		},
	}

	points := []Point{
		{X: 10, Y: 10, Data: "p1"},
		{X: 20, Y: 20, Data: "p2"},
		{X: 30, Y: 30, Data: "p3"},
		{X: 40, Y: 40, Data: "p4"},
		{X: 50, Y: 50, Data: "p5"},
	}

	for _, p := range points {
		qt.Insert(p)
	}

	target := Point{X: 25, Y: 25, Data: nil}

	// Call KNearest multiple times
	result1 := qt.KNearest(target, 3)
	result2 := qt.KNearest(target, 3)

	if len(result1) != len(result2) {
		t.Errorf("Inconsistent results: first call %d, second call %d", len(result1), len(result2))
	}

	// Results should be the same (same order)
	for i := range result1 {
		if result1[i].X != result2[i].X || result1[i].Y != result2[i].Y {
			t.Error("KNearest returned different results on identical calls")
		}
	}
}

// TestKNearestConcurrent tests concurrent k-nearest queries
func TestKNearestConcurrent(t *testing.T) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 1000, Height: 1000},
			Capacity: 10,
		},
	}

	// Pre-populate tree
	for i := 0; i < 200; i++ {
		x := float64(i%20) * 50
		y := float64((i/20)%10) * 100
		qt.Insert(Point{X: x, Y: y, Data: fmt.Sprintf("p%d", i)})
	}

	var wg sync.WaitGroup
	numGoroutines := 10

	for g := 0; g < numGoroutines; g++ {
		wg.Add(1)
		go func(goroutineID int) {
			defer wg.Done()
			for i := 0; i < 20; i++ {
				targetX := float64(goroutineID * 100)
				targetY := float64(i * 50)
				target := Point{X: targetX, Y: targetY, Data: nil}
				results := qt.KNearest(target, 5)

				if len(results) > 5 {
					t.Errorf("KNearest returned more than k results")
				}
			}
		}(g)
	}

	wg.Wait()
}

// BenchmarkKNearest benchmarks k-nearest neighbor search
func BenchmarkKNearest(b *testing.B) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 10,
		},
	}

	// Pre-populate with 5000 points
	for i := 0; i < 5000; i++ {
		x := float64(i%100) * 100
		y := float64((i/100)%50) * 200
		qt.Insert(Point{X: x, Y: y, Data: nil})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetX := float64(i%50) * 200
		targetY := float64((i/50)%50) * 200
		target := Point{X: targetX, Y: targetY, Data: nil}
		_ = qt.KNearest(target, 10)
	}
}

// BenchmarkKNearestSmallK benchmarks k-nearest with small k value
func BenchmarkKNearestSmallK(b *testing.B) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 10,
		},
	}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		x := float64(i%50) * 200
		y := float64((i/50)%20) * 500
		qt.Insert(Point{X: x, Y: y, Data: nil})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		targetX := float64(i%25) * 400
		targetY := float64((i/25)%20) * 500
		target := Point{X: targetX, Y: targetY, Data: nil}
		_ = qt.KNearest(target, 3)
	}
}
