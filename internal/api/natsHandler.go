package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/platform/clients"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type natsHandler struct {
	msgBroker *clients.Client
}

func NewnatsHandler(nc *nats.Conn, js jetstream.Stream) *natsHandler {
	return &natsHandler{
		msgBroker: &clients.Client{NC: nc, JS: js},
	}
}

func (n *natsHandler) dumbHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req RouteRequestDTO
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		log.Printf("Failed to decode request body: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  "Invalid JSON: " + err.Error(),
		})
		return
	}
	if err := req.Validate(); err != nil {
		log.Printf("Request validation failed: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}
	data, _ := json.Marshal(req)

}
