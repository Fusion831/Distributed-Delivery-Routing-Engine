# internal/city - City Grid System

## Overview

The `internal/city` package provides geographic partitioning using grid-based spatial decomposition.

## Components

### Grid (`grid.go`)

**Purpose**: Partition city areas into manageable grid cells for hierarchical dispatch operations.

**Key Concepts**:
- **Grid Cells**: Rectangular regions dividing the city
- **Cell Coordinates**: Map physical coordinates to grid indices
- **Hierarchical Queries**: Support efficient multi-level geographic queries

**Typical Operations**:
- Convert GPS coordinates to grid cell
- Query all deliveries/vehicles in a cell
- Find neighboring cells
- Support grid-level statistics and load balancing

**Use Cases**:
- Dividing city into zones for load balancing
- Grouping nearby deliveries for batch processing
- Regional dispatcher assignment
- Hierarchical routing decisions

## Architecture

The grid system typically works alongside the spatial index:
- **Coarse-grained**: Grid cells organize data into regions
- **Fine-grained**: QuadTree provides efficient queries within cells
- **Combined**: Fast first filter (grid) + precise lookup (tree)

## Future Documentation

- Add grid initialization and configuration examples
- Document cell indexing schemes
- Include performance characteristics
- Add examples of grid-based dispatch strategies
