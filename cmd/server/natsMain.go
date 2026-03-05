package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/api"
)

func natsMain() {
	log.Println("Initializing NATS")
	handler := api.NewnatsHandler("nats://localhost:4222") //TODO: Hardcoding for demo, need to change to env
	defer handler.Close()
	mux := http.NewServeMux()
	mux.HandleFunc("/route", handler.DumbHandle)
	log.Println("HTTP routes")
	srv := &http.Server{
		Addr:         ":8081",
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Println("Server starting on http://localhost:8081")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("Server forced shutdown: %v", err)
	}

	log.Println("Server exited cleanly")
}
