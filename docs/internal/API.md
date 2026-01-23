# REST API Documentation

## Overview

The REST API (`internal/api`) serves as the HTTP "Receptionist" for the Distributed Routing Engine. It wraps the core routing engine, exposing pathfinding capabilities through a clean HTTP interface while managing request/response serialization and error handling.

## Endpoints

### POST /route

Requests an optimal path between two coordinates.

#### Request

```json
{
  "start_x": 5,
  "start_y": 10,
  "end_x": 20,
  "end_y": 35
}
```

**Parameters:**
- `start_x` (int, required): X coordinate of the starting point (≥ 0)
- `start_y` (int, required): Y coordinate of the starting point (≥ 0)
- `end_x` (int, required): X coordinate of the destination (≥ 0)
- `end_y` (int, required): Y coordinate of the destination (≥ 0)

#### Response - Success (200 OK)

```json
{
  "path": [
    {"x": 5, "y": 10},
    {"x": 6, "y": 11},
    {"x": 7, "y": 12},
    {"x": 20, "y": 35}
  ],
  "status": "success"
}
```

#### Response - Error (400 Bad Request)

**Invalid JSON:**
```json
{
  "status": "error",
  "error": "Invalid JSON: unexpected character"
}
```

**Invalid Coordinates:**
```json
{
  "status": "error",
  "error": "coordinates cannot be negative: start=(-5,10) end=(20,35)"
}
```

**Out of Grid Bounds:**
```json
{
  "status": "error",
  "error": "Start coordinates out of grid bounds"
}
```

#### Response - Server Busy (503 Service Unavailable)

When the dispatcher queue is full (backpressure):
```json
{
  "status": "error",
  "error": "Server busy, please retry"
}
```

**Note:** Implement retry logic on the client side with exponential backoff when receiving 503.

#### Response - Pathfinding Failed (500 Internal Server Error)

When no valid path exists:
```json
{
  "status": "error",
  "error": "no path found between start and end nodes"
}
```

## Status Codes

| Code | Meaning | Cause |
|------|---------|-------|
| 200 | Success | Path found and returned |
| 400 | Bad Request | Invalid JSON, negative coordinates, or out-of-bounds coordinates |
| 503 | Service Unavailable | Dispatcher queue is full (backpressure handling) |
| 500 | Internal Error | Pathfinding failed (no path exists between nodes) |

## Request Flow

1. **Validation**: Handler validates JSON structure and non-negative coordinates
2. **Grid Lookup**: Coordinates are converted to grid nodes via `CityGrid.GetNode()`
3. **Bounds Check**: Verify start and end nodes exist within grid bounds
4. **Context Capture**: HTTP request context is captured for cancellation propagation
5. **Queue Submission**: `RouteRequest` is submitted to the `Dispatcher`
6. **Async Processing**: Pathfinding happens asynchronously via worker pool
7. **Result Retrieval**: Handler awaits result on the `resultChan`
8. **Response Formatting**: Path is converted from `model.Node` to `Coordinate` DTOs
9. **Serialization**: Response is JSON-encoded and sent to client

## Design Patterns

### Data Transfer Objects (DTOs)

Request and response types are separate from internal domain models:
- **RouteRequestDTO/RouteResponseDTO**: API contracts for clients
- **model.Node, dispatch.Job**: Internal implementation details (hidden from clients)

### Async Processing

The API uses channels for asynchronous pathfinding:
```go
resultChan := make(chan *dispatch.RouteResult, 1)
routeReq := &dispatch.RouteRequest{
    Start:      startNode,
    End:        endNode,
    Ctx:        ctx,
    ResultChan: resultChan,
}
h.dispatcher.Submit(routeReq)  // Non-blocking queue submission
result := <-resultChan         // Wait for async result
```

### Context Propagation

The HTTP request context is passed through to the pathfinding engine, enabling:
- Client cancellation detection
- Request timeout handling
- Graceful shutdown support

### Backpressure Handling

If the dispatcher queue is full, the handler returns `503 Service Unavailable`:
```go
if err := h.dispatcher.Submit(routeReq); err != nil {
    w.WriteHeader(http.StatusServiceUnavailable)
    // Return error response
}
```

Clients should implement retry logic with exponential backoff.

## Error Handling Strategy

| Scenario | Status | Handling |
|----------|--------|----------|
| Invalid JSON | 400 | Log error, return validation error message |
| Negative coordinates | 400 | Validate in DTO, reject early |
| Out of bounds | 400 | Check grid existence, provide clear message |
| Queue full | 503 | Reject with backpressure signal |
| No path exists | 500 | Return pathfinding error message |

## Integration Example

### Setting up the API

```go
package main

import (
	"net/http"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/api"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/dispatch"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/city"
)

func main() {
	// Initialize components
	dispatcher := dispatch.NewDispatcher(10) // 10 workers
	grid := city.NewCityGrid(100, 100)
	
	// Create API handler
	handler := api.NewHandler(dispatcher, grid)
	
	// Register route
	http.HandleFunc("/route", handler.HandleRoute)
	
	// Start server
	http.ListenAndServe(":8080", nil)
}
```

### cURL Examples

**Successful route request:**
```bash
curl -X POST http://localhost:8080/route \
  -H "Content-Type: application/json" \
  -d '{
    "start_x": 5,
    "start_y": 10,
    "end_x": 20,
    "end_y": 35
  }'
```

**Invalid coordinates (negative values):**
```bash
curl -X POST http://localhost:8080/route \
  -H "Content-Type: application/json" \
  -d '{
    "start_x": -5,
    "start_y": 10,
    "end_x": 20,
    "end_y": 35
  }'
```

**Out of bounds:**
```bash
curl -X POST http://localhost:8080/route \
  -H "Content-Type: application/json" \
  -d '{
    "start_x": 500,
    "start_y": 500,
    "end_x": 600,
    "end_y": 600
  }'
```

**Using jq for pretty-printed response:**
```bash
curl -X POST http://localhost:8080/route \
  -H "Content-Type: application/json" \
  -d '{"start_x": 5, "start_y": 10, "end_x": 20, "end_y": 35}' | jq
```

## Performance Considerations

- **Request Handling**: O(1) - Direct grid lookup
- **Queue Submission**: O(1) - Channel send operation
- **Pathfinding**: O(n log n) - Performed asynchronously by A* algorithm
- **Response Encoding**: O(k) - k = path length

## Future Enhancements

- [ ] Request timeout enforcement
- [ ] Rate limiting per client
- [ ] Batch route requests endpoint
- [ ] WebSocket support for streaming updates
- [ ] Metrics/observability endpoints
- [ ] API versioning
- [ ] Authentication/authorization layer
