package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/api"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/city"
	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/dispatch"
)

func main() {
	// 1. Initialize the routing infrastructure
	log.Println("Initializing routing engine...")

	// Create the city grid (100x100)
	grid := city.NewCityGrid(100, 100)
	log.Println("✓ City grid created (100x100)")

	// Initialize all grid nodes
	grid.InitializeNodes()
	log.Println("✓ Grid nodes initialized")

	// Create the dispatcher with 4 workers and buffer size of 100
	dispatcher := dispatch.NewDispatcher(4, 100, grid)
	log.Println("✓ Dispatcher created (4 workers, buffer=100)")

	// Start the dispatcher worker pool
	go dispatcher.Start()
	log.Println("✓ Dispatcher workers started")

	// 2. Create the API handler
	handler := api.NewHandler(dispatcher, grid)
	log.Println("✓ API handler initialized")

	mux := http.NewServeMux()
	mux.HandleFunc("/route", handler.HandleRoute)
	mux.HandleFunc("/health", handleHealth)
	log.Println("✓ HTTP routes registered")

	// 4. Create HTTP server
	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
		// Add timeouts for production safety
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// 5. Start server in background
	go func() {
		log.Println("Server starting on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// 6. Setup graceful shutdown on signal
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	// Block until signal is received
	<-quit
	log.Println("\n📴 Shutdown signal received, gracefully stopping...")

	// 7. Graceful shutdown with 5-second timeout
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Stop accepting new requests and wait for active ones to finish
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("⚠️  Server forced shutdown: %v", err)
	}

	// 8. Stop the dispatcher gracefully
	if err := dispatcher.GracefulStop(ctx); err != nil {
		log.Printf("⚠️  Dispatcher shutdown timeout: %v", err)
	}

	log.Println("✓ Server exited cleanly")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
