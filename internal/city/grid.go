// Package city provides geographic partitioning using grid-based spatial decomposition.
// The grid system partitions city areas into manageable cells for hierarchical dispatch operations,
// working alongside spatial indices for efficient vehicle and delivery queries.
package city

import (
	"context"
	"sync"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/algo"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// CityGrid partitions a 2D city area into a grid structure for hierarchical dispatch operations.
// Supports efficient neighbor queries and pathfinding using A* algorithm.
// Maintains thread-safety with RWMutex for concurrent access.
type CityGrid struct {
	width, height float64         // Total dimensions of the city grid
	nodes         [][]*model.Node // 2D grid of nodes
	Lock          sync.RWMutex    // Protects concurrent access to grid nodes
}

// NewCityGrid creates and returns a new CityGrid with specified dimensions.
// The grid is initialized but nodes must be populated separately by calling InitializeNodes().
func NewCityGrid(width, height float64) *CityGrid {
	return &CityGrid{
		width:  width,
		height: height,
	}
}

// InitializeNodes populates the grid with nodes for each coordinate.
// Creates width × height nodes, with no obstacles. Must be called before pathfinding.
func (g *CityGrid) InitializeNodes() {
	g.Lock.Lock()
	defer g.Lock.Unlock()

	// Create the 2D node grid
	g.nodes = make([][]*model.Node, int(g.width))
	for x := 0; x < int(g.width); x++ {
		g.nodes[x] = make([]*model.Node, int(g.height))
		for y := 0; y < int(g.height); y++ {
			g.nodes[x][y] = &model.Node{
				X:          x,
				Y:          y,
				Weight:     1,     // Default traversal cost
				IsObstacle: false, // No obstacles by default
			}
		}
	}
}

// GetNode retrieves the node at the given coordinates.
// Returns nil if coordinates are out of bounds.
// Thread-safe with read lock.
func (g *CityGrid) GetNode(x, y int) *model.Node {
	g.Lock.RLock()
	defer g.Lock.RUnlock()

	if x < 0 || x >= len(g.nodes) || y < 0 || y >= len(g.nodes[0]) {
		return nil
	}
	return g.nodes[x][y]
}

// GetNeighbors returns the valid neighboring nodes of a given node.
// Checks all four cardinal directions (up, down, left, right) and filters out obstacles.
// Used by A* algorithm for pathfinding on the grid.
func (g *CityGrid) GetNeighbors(n *model.Node) []*model.Node {
	g.Lock.RLock()
	defer g.Lock.RUnlock()
	var neighbors []*model.Node
	directions := [][2]int{
		{-1, 0},
		{0, 1},
		{0, -1},
		{1, 0},
	}
	for _, direction := range directions {
		newX := n.X + direction[0]
		newY := n.Y + direction[1]
		if newX >= 0 && newX < len(g.nodes) && newY >= 0 && newY < len(g.nodes[0]) {
			neighbor := g.nodes[newX][newY]
			if neighbor != nil && !neighbor.IsObstacle {
				neighbors = append(neighbors, neighbor)
			}
		}

	}
	return neighbors
}

// FindPath computes the shortest path between two nodes on the city grid.
// Uses the A* algorithm via the algo package. Supports context cancellation.
//
// Parameters:
//   - ctx: context for cancellation support
//   - start: the starting node
//   - end: the destination node
//
// Returns the path as a slice of nodes, or an error if no path exists.
func (g *CityGrid) FindPath(ctx context.Context, start, end *model.Node) ([]*model.Node, error) {
	return algo.FindPath(ctx, g, start, end)
}
