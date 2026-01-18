# pkg/algo - Routing Algorithms

## Overview

The `pkg/algo` package provides graph-based and heuristic algorithms for pathfinding and route optimization.

## Components

### A* Pathfinding (`astar.go`)

**Purpose**: Find optimal paths between two points on a graph.

**Algorithm**: A* (A-Star) is a graph traversal algorithm that combines Dijkstra's method with a heuristic estimate of the remaining distance.

**Key Characteristics**:
- Admissible heuristic ensures finding optimal path
- Efficient in practice with good heuristic function
- Suitable for real-time route planning

**Use Cases**:
- Computing shortest route between delivery points
- Finding alternative routes around obstacles
- Real-time navigation queries

### Binary Heap (`heap.go`)

**Purpose**: Efficient priority queue implementation for algorithm operations.

**Data Structure**: Binary heap provides:
- O(1) peek at min/max element
- O(log n) insertion and extraction
- Used by A* algorithm for efficient open set management

## Future Documentation

- Add detailed API reference with examples
- Include algorithm complexity analysis
- Document configuration parameters
- Add usage examples for route optimization
