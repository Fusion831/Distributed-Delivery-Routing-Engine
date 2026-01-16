package dispatch

import (
	"context"
	"testing"
	"time"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/pkg/model"
)

// MockGraph is a test implementation of model.Graph that simulates slow operations
type MockGraph struct {
	sleepDuration time.Duration
	endNode       *model.Node
}

func (mg *MockGraph) GetNeighbors(n *model.Node) []*model.Node {
	time.Sleep(mg.sleepDuration)
	// Return the end node directly so A* terminates quickly
	if mg.endNode != nil {
		return []*model.Node{mg.endNode}
	}
	return []*model.Node{}
}

func (mg *MockGraph) FindPath(ctx context.Context, start, end *model.Node) ([]*model.Node, error) {
	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case <-time.After(mg.sleepDuration):
		// Return a simple path from start to end
		return []*model.Node{start, end}, nil
	}
}

// TestDispatcherBufferOverflow verifies that Submit returns "server busy" when buffer is full
func TestDispatcherBufferOverflow(t *testing.T) {
	// Create mock nodes
	startNode := &model.Node{X: 0, Y: 0}
	endNode := &model.Node{X: 10, Y: 10}

	// Setup: Create Dispatcher with 1 Worker and Buffer Size = 1
	mockGraph := &MockGraph{sleepDuration: 100 * time.Millisecond, endNode: endNode}
	dispatcher := NewDispatcher(1, 1, mockGraph)

	// Start the dispatcher in a goroutine
	go dispatcher.Start()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Job A: Send first request - Worker should grab it immediately
	jobA := &RouteRequest{
		Start:      startNode,
		End:        endNode,
		Ctx:        ctx,
		ResultChan: make(chan *RouteResult, 1),
	}

	errA := dispatcher.Submit(jobA)
	if errA != nil {
		t.Fatalf("Job A submit failed: %v", errA)
	}

	// Give the worker a brief moment to grab Job A from the queue
	time.Sleep(10 * time.Millisecond)

	// Job B: Send second request - Should go into buffer (buffer is now full)
	jobB := &RouteRequest{
		Start:      startNode,
		End:        endNode,
		Ctx:        ctx,
		ResultChan: make(chan *RouteResult, 1),
	}

	errB := dispatcher.Submit(jobB)
	if errB != nil {
		t.Fatalf("Job B submit failed: %v", errB)
	}

	// Job C: Send third request - Buffer is full, should return "server busy" error immediately
	jobC := &RouteRequest{
		Start:      startNode,
		End:        endNode,
		Ctx:        ctx,
		ResultChan: make(chan *RouteResult, 1),
	}

	startTime := time.Now()
	errC := dispatcher.Submit(jobC)
	elapsed := time.Since(startTime)

	// Assert: Job C must return server busy error
	if errC == nil {
		t.Fatal("Job C should have returned an error, but got nil")
	}

	expectedErr := "server busy"
	if errC.Error() != expectedErr {
		t.Errorf("Job C returned wrong error. Expected: %q, Got: %q", expectedErr, errC.Error())
	}

	// Assert: Job C must return instantly (not block for 100ms)
	if elapsed > 50*time.Millisecond {
		t.Errorf("Job C did not return instantly. Took: %v", elapsed)
	}

	// Cleanup: Read results from Job A and B to avoid goroutine leaks
	resultA := <-jobA.ResultChan
	if resultA.err != nil {
		t.Errorf("Job A returned error: %v", resultA.err)
	}

	resultB := <-jobB.ResultChan
	if resultB.err != nil {
		t.Errorf("Job B returned error: %v", resultB.err)
	}

	// Stop the dispatcher gracefully
	dispatcher.Stop()
}
