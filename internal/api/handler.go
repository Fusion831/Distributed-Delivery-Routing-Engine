// Package api implements HTTP request handlers for the routing engine.
package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/city"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/dispatch"
)

// Handler implements the HTTP "Receptionist" that receives route requests,
// delegates to the Dispatcher, and returns results to clients.
// It manages the interaction between the HTTP layer and the core routing engine.
type Handler struct {
	dispatcher *dispatch.Dispatcher // Routes requests to worker pool
	grid       *city.CityGrid       // Provides grid node lookup by coordinates
}

// NewHandler creates a new API handler with the given dispatcher and grid.
// Both dependencies are required for the handler to function.
func NewHandler(d *dispatch.Dispatcher, g *city.CityGrid) *Handler {
	return &Handler{
		dispatcher: d,
		grid:       g,
	}
}

// HandleRoute processes HTTP POST requests for pathfinding.
// It implements http.HandlerFunc signature for use with http.ServeMux.
//
// Request: POST /route with JSON body containing start and end coordinates.
// Response: JSON with path, status, and optional error message.
//
// Error responses:
//   - 400 Bad Request: Invalid JSON or coordinates
//   - 503 Service Unavailable: Dispatcher queue is full (backpressure)
//   - 500 Internal Server Error: Pathfinding failed (no path found)
func (h *Handler) HandleRoute(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Decode and validate request
	var req RouteRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  "Invalid JSON: " + err.Error(),
		})
		return
	}

	if err := req.Validate(); err != nil {
		log.Printf("Request validation failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}

	// Convert DTO coordinates to internal model.Node pointers from the grid
	startNode := h.grid.GetNode(req.StartX, req.StartY)
	if startNode == nil {
		log.Printf("Invalid start coordinates: (%d, %d)", req.StartX, req.StartY)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  "Start coordinates out of grid bounds",
		})
		return
	}

	endNode := h.grid.GetNode(req.EndX, req.EndY)
	if endNode == nil {
		log.Printf("Invalid end coordinates: (%d, %d)", req.EndX, req.EndY)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  "End coordinates out of grid bounds",
		})
		return
	}

	// Capture the HTTP request context for cancellation propagation.
	// If the client closes the connection, r.Context() will be cancelled.
	ctx := r.Context()

	// Create a result channel for this request
	resultChan := make(chan *dispatch.RouteResult, 1)

	// Create and submit the route request to the dispatcher
	routeReq := &dispatch.RouteRequest{
		Start:      startNode,
		End:        endNode,
		Ctx:        ctx,
		ResultChan: resultChan,
	}

	// Submit to dispatcher, handling backpressure
	if err := h.dispatcher.Submit(routeReq); err != nil {
		log.Printf("Dispatcher queue full: %v", err)
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  "Server busy, please retry", //TODO: Implement Retry logic for frontend if server is full
		})
		return
	}

	// Wait for the pathfinding result or context cancellation
	result := <-resultChan

	// Handle pathfinding error
	if result.Err != nil {
		log.Printf("Pathfinding error: %v", result.Err)
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  result.Err.Error(),
		})
		return
	}

	// Convert internal model.Node path to API Coordinate slice
	coordinates := make([]Coordinate, len(result.PathResult))
	for i, node := range result.PathResult {
		coordinates[i] = Coordinate{X: node.X, Y: node.Y}
	}

	// Return success response
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(RouteResponseDTO{
		Path:   coordinates,
		Status: "success",
	})
}
