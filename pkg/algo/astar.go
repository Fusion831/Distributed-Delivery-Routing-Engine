// Package algo provides graph-based pathfinding and routing algorithms.
// It includes A* algorithm implementation combined with binary heap for efficient pathfinding,
// and heuristic-based route optimization for real-time dispatch operations.
package algo

import (
	"container/heap"
	"context"
	"errors"
	"math"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// ManHattanDistance computes the Manhattan distance between two grid nodes.
// This serves as an admissible heuristic for the A* algorithm, ensuring optimal pathfinding.
func ManHattanDistance(a, b *model.Node) float64 {
	return math.Abs(float64(a.X-b.X) + float64(a.Y-b.Y))
}

// FindPath uses the A* algorithm to find the shortest path between two nodes on a graph.
// It combines g-score (cost from start) with h-score (heuristic distance to end) for efficient pathfinding.
// The function respects context cancellation for early termination of long-running searches.
//
// Parameters:
//   - ctx: context for cancellation support
//   - g: the graph to search on
//   - start: the starting node
//   - end: the destination node
//
// Returns the path as a slice of nodes from start to end, or an error if no path exists.
func FindPath(ctx context.Context, g model.Graph, start, end *model.Node) ([]*model.Node, error) {
	// Special case: start is already the destination
	if start == end {
		return []*model.Node{start}, nil
	}

	pq := &PriorityQueue{}
	heap.Init(pq)

	//Tracking Cost through a Map, Cost from start Node -> Current Node
	gScore := make(map[*model.Node]float64)
	gScore[start] = 0

	cameFrom := make(map[*model.Node]*model.Node) //Child-> Parent for Path Reconstruction

	openSetMap := make(map[*model.Node]*Item) //FastLookup for optimization, as to avoid infinite loops
	startItem := &Item{
		Node:     start,
		Priority: ManHattanDistance(start, end),
	}

	heap.Push(pq, startItem)
	openSetMap[start] = startItem
	for pq.Len() > 0 {
		//Context Cancellation Check(if person closes app or website)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		default:
		}
		current := heap.Pop(pq).(*Item).Node
		delete(openSetMap, current) //it has been explored, so remove it from the map

		if current == end {
			return ReConstructPath(cameFrom, current), nil
		}

		for _, neighbor := range g.GetNeighbors(current) {
			tentativeG := gScore[current] + 1.0 + float64(neighbor.Weight)
			if oldG, exists := gScore[neighbor]; !exists || tentativeG < oldG {
				cameFrom[neighbor] = current
				gScore[neighbor] = tentativeG
				fScore := tentativeG + ManHattanDistance(neighbor, end)
				newItem := &Item{Node: neighbor, Priority: fScore}
				heap.Push(pq, newItem)
				openSetMap[neighbor] = newItem
			}
		}
	}
	return nil, errors.New("no path found")
}

func ReConstructPath(cameFrom map[*model.Node]*model.Node, current *model.Node) []*model.Node {
	path := []*model.Node{current}
	for parent, exists := cameFrom[current]; exists; parent, exists = cameFrom[current] {
		path = append([]*model.Node{current}, path...)
		current = parent
	}
	return path
}
