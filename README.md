# distributed-delivery-routing-engine

A routing optimization system built in Go focusing on core spatial indexing, graph algorithms, and dispatch logic.

## Features

- **Spatial Indexing**: QuadTree-based 2D spatial index for efficient location queries (O(log n))
- **Pathfinding**: A* algorithm for optimal route planning
- **Grid Partitioning**: City-based geographic grid for hierarchical dispatch
- **Thread-Safe Operations**: RWMutex-protected concurrent access patterns
- **Comprehensive Testing**: 20+ unit tests covering edge cases and concurrent scenarios

## Installation & Setup

### Prerequisites

- **Go**: 1.18 or higher ([Download Go](https://go.dev/dl/))
- **Git**: For cloning the repository
- **Make** (optional): For simplified build commands

### Clone the Repository

```bash
git clone https://github.com/Fusion831/Distributed-Delivery-Routing-Engine.git
cd Distributed-Delivery-Routing-Engine
```

### Download Dependencies

This project uses only the Go standard library, so no external dependencies need to be installed:

```bash
go mod download
```

### Build the Project

**Build the server:**
```bash
go build -o bin/server cmd/server/main.go
```

**Build the ingestor:**
```bash
go build -o bin/ingestor cmd/ingestor/main.go
```

**Build both:**
```bash
go build -o bin/ ./cmd/...
```

### Verify Installation

Run the test suite to verify everything is set up correctly:

```bash
go test ./...
```

You should see output similar to:
```
ok      github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/spatial      0.456s
ok      github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/algo         0.123s
ok      github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/city    0.234s
ok      github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/dispatch 0.345s
```

## Project Structure

```
pkg/
├── spatial/        # QuadTree spatial indexing
├── algo/           # Pathfinding algorithms (A*, heap)
└── model/          # Core data structures (graph)

internal/
├── dispatch/       # Dispatcher & assignment logic
└── city/           # City grid partitioning

cmd/
├── ingestor/       # Data ingestion entry point
└── server/         # Server/API entry point
```

## Documentation

| Component | Overview | Details |
|-----------|----------|---------|
| **REST API** | HTTP request handling & routing | [docs/internal/API.md](docs/internal/API.md) |
| **Dispatcher** | Assignment logic & worker coordination | [docs/internal/DISPATCH.md](docs/internal/DISPATCH.md) |
| **City Grid** | Geographic partitioning | [docs/internal/CITY.md](docs/internal/CITY.md) |
| **Spatial** | 2D spatial indexing with QuadTree | [docs/pkg/SPATIAL.md](docs/pkg/SPATIAL.md) |
| **Algorithms** | A* pathfinding & utilities | [docs/pkg/ALGO.md](docs/pkg/ALGO.md) |
| **Data Models** | Core graph structures | [docs/pkg/MODEL.md](docs/pkg/MODEL.md) |

## Quick Start

### Running the HTTP Server

Start the routing engine server to accept HTTP requests:

```bash
go run cmd/server/main.go
```

The server listens on `http://localhost:8080` by default.

### HTTP Request Example

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

**Response:**
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

See [docs/internal/API.md](docs/internal/API.md) for complete API documentation with all endpoints, error responses, and cURL examples.

### Spatial Index Example

```go
import "pkg/spatial"

// Create QuadTree
qt := &spatial.QuadTree{
    Root: &spatial.Node{
        Bounds:   spatial.Bounds{X: -90, Y: -180, Width: 180, Height: 360},
        Capacity: 32,
    },
}

// Insert point
driver := spatial.Point{X: 40.7128, Y: -74.0060, Data: "driver_1"}
qt.Insert(driver)

// Search region
zone := spatial.Bounds{X: 40.7, Y: -74.1, Width: 0.1, Height: 0.1}
results := qt.Search(zone)

// Find 5 nearest points
nearest := qt.KNearest(driver, 5)
```

## Testing Guide

### Running All Tests

```bash
go test ./...
```

### Running Tests by Package

```bash
# Test spatial indexing (QuadTree)
go test ./pkg/spatial -v

# Test algorithms (A*, heap)
go test ./pkg/algo -v

# Test dispatcher
go test ./internal/dispatch -v

# Test city grid
go test ./internal/city -v
```

### Running Specific Tests

```bash
# Run a specific test function
go test ./pkg/spatial -run TestInsert -v

# Run tests matching a pattern
go test ./... -run "TestConcurrent" -v
```

### Test Coverage

Generate coverage reports:

```bash
# Overall coverage
go test ./... -cover

# Generate HTML coverage report
go test ./... -coverprofile=coverage.out
go tool cover -html=coverage.out -o coverage.html
```

### Benchmark Tests

Measure performance of critical operations:

```bash
# Run benchmarks
go test ./pkg/spatial -bench=. -benchmem

# Compare benchmark results
go test ./pkg/spatial -bench=. -benchmem -count=5
```

### Test Organization

| Package | Tests | Focus |
|---------|-------|-------|
| `pkg/spatial` | QuadTree insert, search, nearest neighbor | Spatial indexing correctness & performance |
| `pkg/algo` | A* algorithm, heap operations | Pathfinding accuracy |
| `internal/city` | Grid node lookup, partitioning | Geographic partitioning |
| `internal/dispatch` | Task queuing, worker coordination | Concurrent request handling |

### Writing New Tests

Following the project's test structure:

```go
func TestMyFeature(t *testing.T) {
    // Arrange: Set up test data
    testData := setupTestData()
    
    // Act: Perform the operation
    result := myFunction(testData)
    
    // Assert: Verify the result
    if result != expected {
        t.Errorf("expected %v, got %v", expected, result)
    }
}
```

Place tests in `*_test.go` files in the same package as the code being tested.

## Performance

| Operation | Complexity | Notes |
|-----------|-----------|-------|
| Insert | O(log n) | Hierarchical tree insertion |
| Search | O(log n + k) | k = results returned |
| Remove | O(log n) | Efficient leaf deletion |
| KNearest | O(R log n) | R = radius expansions |

## Dependencies

- **Go**: 1.18+
- Built with standard library only

## Future Enhancements

- [ ] Environment variable configuration (PORT, WORKERS, GRID_WIDTH, GRID_HEIGHT)
- [ ] Request timeout enforcement
- [ ] Rate limiting per client
- [ ] Batch route requests endpoint
- [ ] WebSocket support for streaming updates
- [ ] Metrics/observability endpoints
- [ ] API versioning
- [ ] Authentication/authorization layer
