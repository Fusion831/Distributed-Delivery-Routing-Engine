package dispatch

import (
	"context"
	"errors"
	"sync"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/algo"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

type Dispatcher struct {
	jobQueue chan *RouteRequest
	workers  int
	graph    model.Graph
	stopCh   chan struct{}
	wg       sync.WaitGroup
}

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
				err:        err,
			}

			// The Cleanup: Close channel to signal completion
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

// Submit - Non-blocking entry point for route requests
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

// GracefulStop - Stops accepting new jobs and waits for workers to finish with timeout
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
