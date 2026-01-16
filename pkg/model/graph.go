package model

import "context"

type Node struct {
	X, Y       int
	Weight     int //Simulate actual roads/weights
	IsObstacle bool
}

type Graph interface {
	GetNeighbors(n *Node) []*Node
	FindPath(ctx context.Context, start, end *Node) ([]*Node, error)
}
