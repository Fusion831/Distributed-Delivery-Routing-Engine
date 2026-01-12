package algo

import (
	"container/heap"
	"context"
	"errors"
	"math"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

//Heap Functions I can use -> len,less,swap,push,pop,update

func ManHattanDistance(a, b *model.Node) float64 {
	//Will serve as the heuristic for the A* algorithm
	return math.Abs(float64(a.X-b.X) + float64(a.Y-b.Y))
}

func FindPath(ctx context.Context, g model.Graph, start, end *model.Node) ([]*model.Node, error) {
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
