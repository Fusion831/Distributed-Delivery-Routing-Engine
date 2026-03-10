package api

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/Fusion831/Distributed-Delivery-Routing-Engine/internal/platform/clients"
	"github.com/nats-io/nats.go"
)

type natsHandler struct {
	msgBroker *clients.Client
}

// NewnatsHandler creates a handler connected to the given NATS URL.
// It returns an error when the NATS client cannot be created so callers
// can decide whether to retry or fail fast.
func NewnatsHandler(URL string) (*natsHandler, error) {
	client, err := clients.NewClient(URL)
	if err != nil {
		return nil, err
	}
	return &natsHandler{
		msgBroker: client,
	}, nil
}

func (n *natsHandler) Close() {
	if n.msgBroker != nil {
		n.msgBroker.Close()
	}
}

func (n *natsHandler) DumbHandle(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if n.msgBroker == nil || n.msgBroker.NC == nil {
		log.Printf("NATS client not initialized")
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(RouteResponseDTO{Status: "error", Error: "message broker unavailable"})
		return
	}
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
	inbox := nats.NewInbox()
	sub, err := n.msgBroker.NC.SubscribeSync(inbox)
	defer sub.Unsubscribe()
	msg := nats.Msg{Subject: "ORDERS.new", Reply: inbox, Data: data}
	n.msgBroker.NC.PublishMsg(&msg)
	responseMsg, err := sub.NextMsg(time.Second * 5)
	if err != nil {
		log.Printf("Request Timed Out: %v", err)
		w.WriteHeader(http.StatusGatewayTimeout)
		json.NewEncoder(w).Encode(RouteResponseDTO{
			Status: "error",
			Error:  err.Error(),
		})
		return
	}
	w.WriteHeader(http.StatusOK)
	w.Write(responseMsg.Data)
}
