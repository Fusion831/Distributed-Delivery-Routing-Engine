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
	JS jetstream.Stream
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
	js, err := jetstream.New(nc)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s, _ := js.CreateStream(ctx, jetstream.StreamConfig{
		Name:      "ORDERS",
		Subjects:  []string{"ORDERS.*"},
		Storage:   jetstream.FileStorage,
		Retention: jetstream.WorkQueuePolicy,
	})
	if err != nil {
		log.Printf("Note: Stream configuration check: %v", err)
	} else {
		log.Printf("Stream  initialized successfully")
	}

	return &Client{
		NC: nc,
		JS: s,
	}, nil

}

func (c *Client) Close() {
	if c.NC != nil {
		c.NC.Drain()
		c.NC.Close()
	}
}
