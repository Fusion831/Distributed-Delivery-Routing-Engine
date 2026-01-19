// Package algo provides graph-based pathfinding and routing algorithms.
// It includes A* algorithm implementation combined with binary heap for efficient pathfinding,
// and heuristic-based route optimization for real-time dispatch operations.
package algo

import (
	"container/heap"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// Item represents a node in the priority queue with its associated priority score.
// Used for A* algorithm to maintain open set of nodes to be explored.
type Item struct {
	Node     *model.Node // The grid node being evaluated
	Priority float64     // F-score: G-score + H-score heuristic
	Index    int         // Index in the heap (maintained by heap.Interface)
}

// PriorityQueue is a min-heap of Items ordered by priority score.
// Lower priority values are explored first in the A* algorithm.
type PriorityQueue []*Item

// Len returns the number of items in the priority queue.
func (pq PriorityQueue) Len() int { return len(pq) }

// Less compares priorities between two items.
// Returns true if item i has lower priority (should be explored first) than item j.
func (pq PriorityQueue) Less(i, j int) bool {
	// We want Pop to give us the lowest priority, not highest, so we use less than here.
	return pq[i].Priority < pq[j].Priority
}

// Swap exchanges two items in the heap.
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

// Push adds a new item to the priority queue.
func (pq *PriorityQueue) Push(x any) {
	n := len(*pq)
	item := x.(*Item)
	item.Index = n
	*pq = append(*pq, item)
}

// Pop removes and returns the item with the lowest priority from the queue.
func (pq *PriorityQueue) Pop() any {
	old := *pq
	n := len(old)
	item := old[n-1]
	old[n-1] = nil  // don't stop the GC from reclaiming the item eventually
	item.Index = -1 // for safety
	*pq = old[0 : n-1]
	return item
}

func (pq *PriorityQueue) update(item *Item, value string, priority float64) {
	item.Priority = float64(priority)
	heap.Fix(pq, item.Index)
}
