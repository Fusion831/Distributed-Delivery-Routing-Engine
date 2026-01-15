package dispatch

import (
	"context"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

type RouteResult struct {
	PathResult []*model.Node
	err        error
}

type RouteRequest struct {
	Start, End *model.Node
	Ctx        context.Context
	ResultChan chan *RouteResult
}
