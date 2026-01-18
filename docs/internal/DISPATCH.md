# internal/dispatch - Dispatcher & Assignment Logic

## Overview

The `internal/dispatch` package implements core dispatch decision-making and vehicle-to-delivery assignment logic.

## Components

### Dispatcher (`dispatcher.go`)

**Purpose**: Main dispatcher service that manages vehicle assignment and delivery routing.

**Key Responsibilities**:
- Receive incoming delivery requests
- Query available vehicles from spatial index
- Execute assignment algorithms to match vehicles to deliveries
- Manage vehicle state (available, assigned, unavailable)
- Track active assignments
- Handle real-time updates to vehicle locations and capacity

**Typical Workflow**:
1. Delivery request arrives
2. Query spatial index for nearby available vehicles
3. Filter candidates by capacity, constraints
4. Run assignment algorithm (greedy, optimal, heuristic)
5. Update vehicle state and create assignment
6. Return assignment to requesting system

### Type Definitions (`types.go`)

**Purpose**: Define data structures for dispatch operations.

**Typical Types**:
- `Delivery`: Represents a delivery order with location and requirements
- `Vehicle`: Represents an available vehicle with capacity and location
- `Assignment`: Represents dispatch decision linking delivery to vehicle
- `DispatchResult`: Return value indicating success/failure of assignment

**Constraints**:
- Vehicle capacity (weight, volume, item count)
- Delivery time windows
- Vehicle availability
- Geographic constraints

### Tests (`dispatcher_test.go`)

**Purpose**: Validate dispatcher logic and assignment algorithms.

**Test Coverage**:
- Basic assignment scenarios
- Constraint validation
- Concurrent dispatch operations
- Edge cases and error handling

## Assignment Strategies

Different algorithms may be used depending on optimization goals:

- **Greedy**: Assign to nearest available vehicle (fast, suboptimal)
- **Optimal**: Find best vehicle considering all factors (slow, optimal)
- **Heuristic**: Balance speed and quality with approximation algorithms
- **Load Balanced**: Distribute work evenly across vehicles

## Concurrency Considerations

The dispatcher handles concurrent operations:
- Multiple delivery requests arriving simultaneously
- Real-time vehicle updates
- Concurrent state modifications
- Must maintain consistency with spatial index

## Future Documentation

- Add detailed assignment algorithm documentation
- Include performance optimization strategies
- Document constraint types and handling
- Add examples of custom assignment strategies
- Include metrics and monitoring details
