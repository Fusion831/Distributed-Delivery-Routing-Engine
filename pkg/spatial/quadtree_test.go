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

// BenchmarkSearch benchmarks search performance
func BenchmarkSearch(b *testing.B) {
	qt := &QuadTree{
		Root: &Node{
			Bounds:   Bounds{X: 0, Y: 0, Width: 10000, Height: 10000},
			Capacity: 10,
		},
	}

	// Pre-populate
	for i := 0; i < 1000; i++ {
		x := float64(i%100) * 100
		y := float64((i/100)%100) * 100
		qt.Insert(Point{X: x, Y: y, Data: nil})
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		searchArea := Bounds{X: float64(i%50) * 100, Y: float64((i/50)%50) * 100, Width: 500, Height: 500}
		_ = qt.Search(searchArea)
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
