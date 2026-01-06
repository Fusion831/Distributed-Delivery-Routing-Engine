// Current Problem: Testing Systems one by one, and stop focusing on the big picture
package city

import (
	"sync"
	"testing"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// This line fails to compile if *CityGrid doesn't satisfy model.Graph
var _ model.Graph = (*CityGrid)(nil)

// Instead of going straight to the 100x100 grid, I should test it in smaller parts to make sure it is working on all edge cases
func create3x3Grid() *CityGrid {
	grid := NewCityGrid(3, 3)
	grid.nodes = make([][]*model.Node, 3)
	for i := 0; i < 3; i++ {
		grid.nodes[i] = make([]*model.Node, 3)
		for j := 0; j < 3; j++ {
			grid.nodes[i][j] = &model.Node{
				X:          i,
				Y:          j,
				IsObstacle: false,
			}
		}
	}
	return grid
}

// Test B: Table-Driven Test for Logic Verification
func TestGetNeighborsTableDriven(t *testing.T) {
	tests := []struct {
		name          string
		targetX       int
		targetY       int
		obstacles     [][2]int // coordinates to mark as obstacles
		expectedCount int
		description   string
	}{
		{
			name:          "Center node",
			targetX:       1,
			targetY:       1,
			obstacles:     nil,
			expectedCount: 4,
			description:   "Center (1,1) should have 4 neighbors: up, down, left, right",
		},
		{
			name:          "Corner (0,0)",
			targetX:       0,
			targetY:       0,
			obstacles:     nil,
			expectedCount: 2,
			description:   "Corner (0,0) should have 2 neighbors: right (0,1), down (1,0)",
		},
		{
			name:          "Obstacle blocking neighbor",
			targetX:       0,
			targetY:       0,
			obstacles:     [][2]int{{1, 0}}, // block down
			expectedCount: 1,
			description:   "Corner (0,0) with obstacle at (1,0) should have 1 neighbor: right (0,1)",
		},
		{
			name:          "Multiple obstacles",
			targetX:       1,
			targetY:       1,
			obstacles:     [][2]int{{0, 1}, {2, 1}}, // block up and down
			expectedCount: 2,
			description:   "Center (1,1) with obstacles up and down should have 2 neighbors: left, right",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			grid := create3x3Grid()

			// Set obstacles
			for _, obs := range tt.obstacles {
				grid.nodes[obs[0]][obs[1]].IsObstacle = true
			}

			targetNode := grid.nodes[tt.targetX][tt.targetY]
			neighbors := grid.GetNeighbors(targetNode)

			if len(neighbors) != tt.expectedCount {
				t.Errorf("%s: expected %d neighbors, got %d. %s",
					tt.name, tt.expectedCount, len(neighbors), tt.description)
			}
		})
	}
}

// Test C: Race Detector Test (Concurrent Read/Write)
func TestGetNeighborsConcurrentSafety(t *testing.T) {
	grid := create3x3Grid()
	var wg sync.WaitGroup
	const iterations = 100

	// Writer goroutine: toggles obstacle at (1,1) 100 times
	wg.Add(1)
	go func() {
		defer wg.Done()
		targetNode := grid.nodes[1][1]
		for i := 0; i < iterations; i++ {
			grid.Lock.Lock()
			targetNode.IsObstacle = (i % 2) == 0 // toggle on/off
			grid.Lock.Unlock()
		}
	}()

	// Reader goroutines: call GetNeighbors 100 times each (10 concurrent readers)
	for r := 0; r < 10; r++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for i := 0; i < iterations; i++ {
				centerNode := grid.nodes[1][1]
				_ = grid.GetNeighbors(centerNode) // discard result, just check for races
			}
		}()
	}

	wg.Wait()

}
