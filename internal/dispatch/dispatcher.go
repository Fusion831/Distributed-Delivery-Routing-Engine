package dispatch

import (
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
