# pkg/spatial - QuadTree Spatial Indexing

## What is a QuadTree?

A QuadTree is a hierarchical spatial data structure that recursively partitions a 2D space into four equal quadrants. It's particularly useful for:

- **Efficient spatial queries**: Finding all points within a rectangular region in O(log n) time on average
- **Nearest neighbor searches**: Quickly locating k-nearest points to a target location
- **Spatial clustering**: Organizing geospatial data for fast retrieval
- **Low latency operations**: Essential for real-time dispatch systems where sub-millisecond lookups are critical

**How it works**: The QuadTree starts with a root node covering the entire bounded space. When a node's capacity is exceeded, it subdivides into 4 children representing:
- **NW (Northwest)**: Top-left quadrant
- **NE (Northeast)**: Top-right quadrant  
- **SW (Southwest)**: Bottom-left quadrant
- **SE (Southeast)**: Bottom-right quadrant

Points are recursively placed into appropriate child nodes based on their coordinates, creating a balanced tree structure that minimizes search depth.

---

## Core Data Structures

### `Point`
Represents a 2D coordinate with optional associated data:
```go
type Point struct {
    X, Y float64      // Coordinates
    Data interface{}  // Arbitrary data (e.g., vehicle ID, delivery ID)
}
```

### `Bounds`
Defines a rectangular region:
```go
type Bounds struct {
    X      float64  // Top-left X coordinate
    Y      float64  // Top-left Y coordinate
    Width  float64
    Height float64
}
```

### `Node`
Individual node in the QuadTree structure:
```go
type Node struct {
    Bounds   Bounds      // Region covered by this node
    Points   []Point     // Points stored in this node (when leaf)
    Capacity int         // Max points before subdivision
    Children [4]*Node    // Child nodes (when subdivided)
}
```

### `QuadTree`
Thread-safe wrapper for the entire tree:
```go
type QuadTree struct {
    Root *Node
    Lock sync.RWMutex  // Protects concurrent access
}
```

---

## Public API Functions

### `Insert(point Point) bool`
Inserts a new point into the QuadTree.
- **Thread-safe**: Uses write lock
- **Returns**: `true` if insertion succeeded, `false` if point is out of bounds
- **Time complexity**: O(log n) average case
- **Example**: Adding a new driver location to the spatial index

```go
driver := Point{X: 40.7128, Y: -74.0060, Data: "driver_123"}
success := qt.Insert(driver)
```

---

### `Search(area Bounds) []Point`
Finds all points within a rectangular region.
- **Thread-safe**: Uses read lock (multiple concurrent searches allowed)
- **Returns**: Slice of all points intersecting the search bounds
- **Time complexity**: O(log n + k) where k is the number of results
- **Example**: Finding all drivers in a geographic zone

```go
zone := Bounds{X: 40.7, Y: -74.1, Width: 0.1, Height: 0.1}
drivers := qt.Search(zone)
```

---

### `Remove(point Point) bool`
Removes a point from the QuadTree.
- **Thread-safe**: Uses write lock
- **Returns**: `true` if point existed and was removed, `false` otherwise
- **Time complexity**: O(log n) average case
- **Example**: Removing a driver who has gone offline

```go
removed := qt.Remove(driverPoint)
```

---

### `Update(oldPoint, newPoint Point) bool`
Atomically updates a point's location (remove old, insert new).
- **Thread-safe**: Uses single write lock for atomicity
- **Validation**: Ensures new point is within bounds before committing
- **Rollback**: Re-inserts old point if new insertion fails
- **Returns**: `true` if update succeeded
- **Time complexity**: O(log n) average case
- **Example**: Updating driver position as they move

```go
oldLocation := Point{X: 40.7, Y: -74.0, Data: "driver_123"}
newLocation := Point{X: 40.71, Y: -74.01, Data: "driver_123"}
success := qt.Update(oldLocation, newLocation)
```

---

### `KNearest(target Point, k int) []Point`
Finds the k nearest points to a target location.
- **Thread-safe**: Uses read lock
- **Algorithm**: 
  - Starts with small search radius (10.0 units)
  - Progressively doubles radius until ≥k points found
  - Caps max results at 100,000 to prevent memory overflow
  - Sorts results by Euclidean distance
- **Returns**: Ordered slice of k nearest points (or fewer if fewer than k exist)
- **Time complexity**: O(R log n + k log k) where R is the search radius expansion factor
- **Example**: Finding 5 nearest available drivers to a delivery request

```go
delivery := Point{X: 40.72, Y: -74.02, Data: "delivery_456"}
nearestDrivers := qt.KNearest(delivery, 5)
```

---

## Internal Helper Functions

### `(b Bounds) Intersects(other Bounds) bool`
Checks if two rectangular regions overlap or touch using Axis-Aligned Bounding Box (AABB) logic.

### `(b Bounds) Contains(point Point) bool`
Tests if a point lies within the bounds (inclusive of edges).

### `(n *Node) SubDivide()`
Splits a node into 4 child nodes when capacity is exceeded.
- Redistributes existing points to appropriate children
- Clears parent's point list after distribution

### `(n *Node) InsertNode(point Point) bool`
Recursive insertion into the tree structure (internal use).

### `(n *Node) SearchTree(searchArea Bounds, resultPoints *[]Point)`
Recursive search through the tree (internal use).

### `(n *Node) RemoveNode(point Point) bool`
Recursive removal from the tree (internal use).
- Uses swap-and-pop strategy for efficient deletion

### `Distance(p1, p2 Point) float64`
Calculates Euclidean distance between two points: $\sqrt{(x_2-x_1)^2 + (y_2-y_1)^2}$

### `sortByDistance(points []PointWithDistance)`
Insertion sort used to order results by distance in `KNearest`.

---

## Test Suite

The test suite in `quadtree_test.go` provides comprehensive coverage:

### Bounds Tests
- **TestBoundsContains**: Validates point-in-bounds detection, including edge cases (points on boundaries, negative coordinates)
- **TestBoundsIntersects**: Tests rectangular intersection detection with overlapping, touching, and separated regions

### Node Structure Tests
- **TestNodeSubDivide**: Verifies correct creation of 4 children with proper bounds
- **TestNodeSubDivideWithPoints**: Ensures points are correctly distributed to children during subdivision

### Basic Operations Tests
- **TestQuadTreeInitialization**: Validates QuadTree creation
- **TestQuadTreeInsertSingle**: Single point insertion
- **TestQuadTreeInsertMultiple**: Multiple sequential insertions
- **TestQuadTreeInsertExceedsCapacity**: Automatic tree subdivision when capacity exceeded
- **TestQuadTreeInsertOutOfBounds**: Rejection of out-of-bounds points

### Search Tests
- **TestQuadTreeSearchBasic**: Finding points in a region
- **TestQuadTreeSearchNoResults**: Handling empty result sets
- **TestQuadTreeSearchAll**: Retrieving all points in tree
- **TestQuadTreeSearchAfterSubdivision**: Correctness of regional searches after tree restructuring
- **TestQuadTreeSearchPrecision**: Float precision handling
- **TestQuadTreeSearchWithData**: Verification that associated data is preserved

### Removal Tests
- **TestQuadTreeRemoveBasic**: Removing existing points
- **TestQuadTreeRemoveNonexistent**: Handling removal of non-existent points

### Advanced Tests
- **TestQuadTreeDeepNesting**: Behavior with deeply nested tree structures
- **TestQuadTreeNegativeCoordinates**: Correctness with negative coordinate spaces
- **TestQuadTreeConcurrentInsert**: 10 goroutines × 50 points each without race conditions
- **TestQuadTreeConcurrentSearchInsert**: Concurrent read/write operations under contention

---

## Thread Safety & Concurrency

The QuadTree uses a single `sync.RWMutex` for all synchronization:

- **Read operations** (Search, KNearest): Acquire read lock - multiple readers can operate concurrently
- **Write operations** (Insert, Remove, Update): Acquire write lock - exclusive access

### Current Limitations & Future Improvements

- **Single global lock**: All write operations contend for a single lock. Under high insertion/update rates, this becomes a bottleneck
  - **Granular locking**: Future optimization could use per-node locks, allowing independent subtrees to be modified concurrently. Trade-offs: higher memory overhead, increased complexity, deadlock risk
  
- **Lock-free data structures**: Alternative approach using atomic operations and CAS (Compare-And-Swap) could eliminate locks entirely for some operations
  
- **Read-Write separation**: Could maintain separate read-optimized and write-optimized indices with eventual consistency

---

## Performance Characteristics

| Operation | Average | Worst Case | Notes |
|-----------|---------|-----------|-------|
| Insert | O(log n) | O(n) | Worst case: all points in one quadrant |
| Search | O(log n + k) | O(n) | k = number of results |
| Remove | O(log n) | O(n) | Same as insert |
| Update | O(log n) | O(n) | Atomic remove + insert |
| KNearest | O(R log n + k log k) | O(n log n) | R = radius expansion iterations |

---

## Usage Example

```go
// Create a QuadTree for a 180x360 degree lat/lon region
qt := &QuadTree{
    Root: &Node{
        Bounds:   Bounds{X: -90, Y: -180, Width: 180, Height: 360},
        Capacity: 32,  // Subdivide when >32 points in a region
    },
}

// Add drivers
driver1 := Point{X: 40.7128, Y: -74.0060, Data: "driver_1"}
driver2 := Point{X: 34.0522, Y: -118.2437, Data: "driver_2"}
qt.Insert(driver1)
qt.Insert(driver2)

// Find drivers in Manhattan area
manhattanZone := Bounds{X: 40.7, Y: -74.1, Width: 0.1, Height: 0.1}
nearby := qt.Search(manhattanZone)

// Find 5 closest drivers to a delivery
delivery := Point{X: 40.7200, Y: -74.0050, Data: "delivery_1"}
closest := qt.KNearest(delivery, 5)

// Update driver position
newLocation := Point{X: 40.7150, Y: -74.0030, Data: "driver_1"}
qt.Update(driver1, newLocation)
```
