package model

type Node struct {
	X, Y       int
	IsObstacle bool
}

type Graph interface {
	GetNeighbors(n *Node) []*Node
}
