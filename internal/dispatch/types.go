// Package dispatch implements core dispatch decision-making and vehicle-to-delivery assignment logic.
// It manages vehicle assignment, delivery routing, and real-time state updates for the routing engine.
package dispatch

import (
	"context"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// RouteResult represents the result of a pathfinding operation.
// Contains either the computed path or an error if pathfinding failed.
type RouteResult struct {
	PathResult []*model.Node // The sequence of nodes forming the shortest path
	err        error         // Error if pathfinding failed (e.g., no path found)
}

// RouteRequest represents a request to find a path between two nodes.
// Uses channels for asynchronous communication in the dispatcher's worker pool.
type RouteRequest struct {
	Start      *model.Node       // The starting node for pathfinding
	End        *model.Node       // The destination node for pathfinding
	Ctx        context.Context   // Context supporting cancellation
	ResultChan chan *RouteResult // Channel to receive the pathfinding result
}
