# distributed-delivery-routing-engine

A routing optimization system built in Go focusing on core spatial indexing, graph algorithms, and dispatch logic.

## Features

- **Spatial Indexing**: QuadTree-based 2D spatial index for efficient location queries (O(log n))
- **Pathfinding**: A* algorithm for optimal route planning
- **Grid Partitioning**: City-based geographic grid for hierarchical dispatch
- **Thread-Safe Operations**: RWMutex-protected concurrent access patterns
- **Comprehensive Testing**: 20+ unit tests covering edge cases and concurrent scenarios

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
| **Spatial** | 2D spatial indexing with QuadTree | [docs/pkg/SPATIAL.md](docs/pkg/SPATIAL.md) |
| **Algorithms** | A* pathfinding & utilities | [docs/pkg/ALGO.md](docs/pkg/ALGO.md) |
| **Data Models** | Core graph structures | [docs/pkg/MODEL.md](docs/pkg/MODEL.md) |
| **City Grid** | Geographic partitioning | [docs/internal/CITY.md](docs/internal/CITY.md) |
| **Dispatcher** | Assignment logic | [docs/internal/DISPATCH.md](docs/internal/DISPATCH.md) |

## Quick Start

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

## Running Tests

```bash
go test ./...
```

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

## Future Documentation

As components are developed or expanded, detailed documentation will be added to:
- `docs/pkg/` - Package-level deep dives
- `docs/internal/` - Internal component guides
- Package godoc comments in source files
