package dispatch

import (
	"errors"
	"sync"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/algo"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

type Dispatcher struct {
	jobQueue chan *RouteRequest
	workers  int
	graph    model.Graph
}

func NewDispatcher(workerCount, bufferSize int, graph model.Graph) *Dispatcher {
	return &Dispatcher{
		jobQueue: make(chan *RouteRequest, bufferSize),
		workers:  workerCount,
		graph:    graph,
	}
}

func (d Dispatcher) workerLoop() {
	for job := range d.jobQueue {
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

func (d Dispatcher) Start() {
	var wg sync.WaitGroup
	for i := 0; i < d.workers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			d.workerLoop()
		}()
	}
	wg.Wait()
}

// Submit - Non-blocking entry point for route requests
func (d Dispatcher) Submit(req *RouteRequest) error {
	select {
	case d.jobQueue <- req:
		return nil
	default:
		return errors.New("server busy")
	}
}
