package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"
)

type RouteRequest struct {
	OriginLat      float64 `json:"origin_lat"`
	OriginLon      float64 `json:"origin_lon"`
	DestinationLat float64 `json:"destination_lat"`
	DestinationLon float64 `json:"destination_lon"`
	VehicleType    string  `json:"vehicle_type"`
}

type RouteResponse struct {
	Route    [][2]float64 `json:"route"`
	Distance float64      `json:"distance"`
	Duration int          `json:"duration"`
	Status   string       `json:"status"`
}

func handleRoute(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req RouteRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Simulate routing calculation (~1ms)
	time.Sleep(time.Millisecond)

	// Calculate mock distance
	dlat := req.DestinationLat - req.OriginLat
	dlon := req.DestinationLon - req.OriginLon
	distance := (dlat*dlat + dlon*dlon) * 111 // rough km estimate

	response := RouteResponse{
		Route: [][2]float64{
			{req.OriginLat, req.OriginLon},
			{req.DestinationLat, req.DestinationLon},
		},
		Distance: distance,
		Duration: int(distance * 60), // minutes
		Status:   "success",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	http.HandleFunc("/route", handleRoute)
	http.HandleFunc("/health", handleHealth)

	addr := fmt.Sprintf(":%s", port)
	log.Printf("Test server starting on http://localhost%s", addr)
	if err := http.ListenAndServe(addr, nil); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
