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

			if job.Ctx.Err() != nil {
				close(job.ResultChan)
				continue
			}

			pathResult, err := algo.FindPath(job.Ctx, d.graph, job.Start, job.End)

			job.ResultChan <- &RouteResult{
				PathResult: pathResult,
				Err:        err,
			}

			close(job.ResultChan)
		}
	}
}

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

func (d *Dispatcher) Stop() {
	close(d.stopCh)
}

// GracefulStop closes the job queue and waits for all workers to finish processing.
// Respects the provided context for timeout. If timeout is exceeded, forces shutdown.
// Returns an error if the context is cancelled before workers finish.
func (d *Dispatcher) GracefulStop(ctx context.Context) error {

	close(d.jobQueue)

	done := make(chan struct{})
	go func() {
		d.wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		return nil
	case <-ctx.Done():

		d.Stop()
		return ctx.Err()
	}
}
