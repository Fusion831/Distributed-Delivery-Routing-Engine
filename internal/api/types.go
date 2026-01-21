// Package api defines the HTTP API contract and request/response types.
// All structs here are Data Transfer Objects (DTOs) for JSON serialization,
// keeping internal domain models (model.Node, dispatch.Job) hidden from clients.
package api

import "fmt"

// RouteRequestDTO represents a client request for pathfinding between two coordinates.
type RouteRequestDTO struct {
	StartX int `json:"start_x"` // X coordinate of the starting point
	StartY int `json:"start_y"` // Y coordinate of the starting point
	EndX   int `json:"end_x"`   // X coordinate of the destination
	EndY   int `json:"end_y"`   // Y coordinate of the destination
}

// Validate checks that all coordinates are non-negative.
// Returns an error if validation fails, suitable for returning as 400 Bad Request.
func (r *RouteRequestDTO) Validate() error {
	if r.StartX < 0 || r.StartY < 0 || r.EndX < 0 || r.EndY < 0 {
		return fmt.Errorf("coordinates cannot be negative: start=(%d,%d) end=(%d,%d)",
			r.StartX, r.StartY, r.EndX, r.EndY)
	}
	return nil
}

// Coordinate represents a single point in 2D space for the API response.
type Coordinate struct {
	X int `json:"x"` // X coordinate
	Y int `json:"y"` // Y coordinate
}

// RouteResponseDTO represents the response from a pathfinding request.
type RouteResponseDTO struct {
	Path   []Coordinate `json:"path"`            // Sequence of coordinates forming the path
	Status string       `json:"status"`          // Status of the request (e.g., "success", "error")
	Error  string       `json:"error,omitempty"` // Error message (omitted from JSON if empty)
}
