// Package dispatch implements core dispatch decision-making and vehicle-to-delivery assignment logic.
// It manages vehicle assignment, delivery routing, and real-time state updates for the routing engine.
package dispatch

import (
	"context"
	"errors"
	"sync"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/algo"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// Dispatcher manages concurrent route requests using a worker pool pattern.
// It processes incoming pathfinding requests and distributes them to available workers
// for A* pathfinding operations on the provided graph.
type Dispatcher struct {
	jobQueue chan *RouteRequest // Channel for incoming route requests
	workers  int                // Number of worker goroutines
	graph    model.Graph        // The graph on which to compute paths
	stopCh   chan struct{}      // Signal channel for stopping workers
	wg       sync.WaitGroup     // Synchronizes worker goroutines
}

// NewDispatcher creates a new Dispatcher with specified worker count, buffer size, and graph.
//
// Parameters:
//   - workerCount: number of worker goroutines for concurrent pathfinding
//   - bufferSize: capacity of the job queue channel
//   - graph: the routing graph for pathfinding operations
//
// Returns a new Dispatcher instance ready to start.
func NewDispatcher(workerCount, bufferSize int, graph model.Graph) *Dispatcher {
	return &Dispatcher{
		jobQueue: make(chan *RouteRequest, bufferSize),
		workers:  workerCount,
		graph:    graph,
		stopCh:   make(chan struct{}),
	}
}

func (d *Dispatcher) workerLoop() {
	defer d.wg.Done()
	for {
		select {
		case <-d.stopCh:
			return
		case job := <-d.jobQueue:
			if job == nil {
				return
			}
			// Stale Check: Skip if context is already cancelled
			if job.Ctx.Err() != nil {
				close(job.ResultChan)
				continue
			}

			pathResult, err := algo.FindPath(job.Ctx, d.graph, job.Start, job.End)

			// The Response: Send result into the channel
			job.ResultChan <- &RouteResult{
				PathResult: pathResult,
				Err:        err,
			}

			// The Cleanup: Close channel to signal completion
			close(job.ResultChan)
		}
	}
}

// Start begins processing route requests with the configured number of worker goroutines.
// This is a blocking call that waits for all workers to finish.
func (d *Dispatcher) Start() {
	for i := 0; i < d.workers; i++ {
		d.wg.Add(1)
		go d.workerLoop()
	}
	d.wg.Wait()
}

// Submit enqueues a route request for processing by an available worker.
// Returns nil if successfully queued, or an error if the queue is full.
// This is a non-blocking operation.
func (d *Dispatcher) Submit(req *RouteRequest) error {
	select {
	case d.jobQueue <- req:
		return nil
	default:
		return errors.New("server busy")
	}
}

// Stop immediately terminates all workers without waiting for pending requests.
func (d *Dispatcher) Stop() {
	close(d.stopCh)
}

// GracefulStop closes the job queue and waits for all workers to finish processing.
// Respects the provided context for timeout. If timeout is exceeded, forces shutdown.
// Returns an error if the context is cancelled before workers finish.
func (d *Dispatcher) GracefulStop(ctx context.Context) error {

	close(d.jobQueue)

	// Wait for workers to finish or timeout
	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():
		// Timeout exceeded, force stop
		d.Stop()
		return ctx.Err()
	}
}
