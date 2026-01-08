package model

type Node struct {
	X, Y       int
	weight     int //Simulate actual roads/weights
	IsObstacle bool
}

type Graph interface {
	GetNeighbors(n *Node) []*Node
}
