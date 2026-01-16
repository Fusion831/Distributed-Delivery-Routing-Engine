package city

import (
	"context"
	"sync"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/algo"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

type CityGrid struct {
	width, height float64
	nodes         [][]*model.Node
	Lock          sync.RWMutex
}

func NewCityGrid(width, height float64) *CityGrid {
	return &CityGrid{
		width:  width,
		height: height,
	}
}

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

func (g *CityGrid) FindPath(ctx context.Context, start, end *model.Node) ([]*model.Node, error) {
	return algo.FindPath(ctx, g, start, end)
}
