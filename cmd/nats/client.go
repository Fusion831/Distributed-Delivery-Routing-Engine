package nats

import (
	"context"
	"log"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Client struct {
	NC *nats.Conn
	JS *jetstream.JetStream
}

func NewClient(url string) (*Client, error) {
	var nc *nats.Conn
	var err error
	for i := 0; i < 5; i++ {
		nc, err = nats.Connect(url, nats.Name("VRP ENGINE"))
		if err == nil {
			break
		}
		log.Printf("retrying in %d Seconds", i+1)
		time.Sleep(time.Duration(i+1) * time.Second)

	}
	if err != nil {
		return nil, err
	}
	js, _ := jetstream.New(nc)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
}
