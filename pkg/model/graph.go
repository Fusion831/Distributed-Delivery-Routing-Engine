// Package model defines core data structures used throughout the routing engine.
// It includes node representations, graph interfaces, and delivery/vehicle models
// for pathfinding and dispatch operations.
package model

import "context"

// Node represents a location in the grid-based routing network.
// Nodes are connected to form a graph used for pathfinding and routing.
type Node struct {
	X          int  // X coordinate in the grid
	Y          int  // Y coordinate in the grid
	Weight     int  // Weight representing traversal cost (distance, time, or terrain difficulty)
	IsObstacle bool // Whether this node is an obstacle and cannot be traversed
}

// Graph defines the interface for routing graphs.
// Implementations provide neighbor lookup and pathfinding capabilities.
type Graph interface {
	// GetNeighbors returns adjacent nodes reachable from the given node.
	GetNeighbors(n *Node) []*Node
	// FindPath finds the shortest path from start to end node using context for cancellation.
	FindPath(ctx context.Context, start, end *Node) ([]*Node, error)
}
