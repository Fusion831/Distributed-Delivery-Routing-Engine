# pkg/model - Data Models

## Overview

The `pkg/model` package defines core data structures used throughout the routing engine.

## Components

### Graph (`graph.go`)

**Purpose**: Represent the delivery/road network as a directed graph.

**Key Concepts**:
- **Nodes**: Represent locations (intersections, delivery points, hubs)
- **Edges**: Represent connections between locations (roads, possible routes)
- **Weights**: Edge weights represent distance, time, or cost

**Typical Structure**:
- Graph nodes maintain adjacency lists of connected nodes
- Edges store routing information (distance, estimated time, capacity)
- Support for directed and weighted edges

**Use Cases**:
- Representing city road networks
- Modeling delivery zone connectivity
- Supporting pathfinding algorithms (A*, Dijkstra)

## Future Documentation

- Add detailed node and edge structure documentation
- Include graph construction examples
- Document graph modification operations
- Add performance characteristics for graph operations
