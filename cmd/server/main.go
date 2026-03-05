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

func shitmain() {

	log.Println("Initializing routing engine...")

	grid := city.NewCityGrid(100, 100)
	log.Println(" City grid created (100x100)")

	grid.InitializeNodes()
	log.Println(" Grid nodes initialized")

	dispatcher := dispatch.NewDispatcher(4, 100, grid)
	log.Println("Dispatcher created (4 workers, buffer=100)")

	go dispatcher.Start()
	log.Println("Dispatcher workers started")

	handler := api.NewHandler(dispatcher, grid)
	log.Println("API handler initialized")

	mux := http.NewServeMux()
	mux.HandleFunc("/route", handler.HandleRoute)
	mux.HandleFunc("/health", handleHealth)
	log.Println(" HTTP routes registered")

	srv := &http.Server{
		Addr:         ":8080",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Println("Server starting on http://localhost:8080")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit
	log.Println("\n📴 Shutdown signal received, gracefully stopping...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
	}

	// 8. Stop the dispatcher gracefully
	if err := dispatcher.GracefulStop(ctx); err != nil {
		log.Printf("Dispatcher shutdown timeout: %v", err)
	}

	log.Println("✓ Server exited cleanly")
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"healthy"}`))
}
