package model

import (
	"fmt"
	"math"
)

type Node struct {
	X,Y int
	isObstacle bool
}

type Graph interface {
	GetNeighbors(n *Node) []*Node
}