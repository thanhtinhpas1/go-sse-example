package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/adaptor"
	uuid "github.com/satori/go.uuid"
)

type Client struct {
	name   string
	events chan *OrderStatus
}

type OrderStatus struct {
	AppTransId string
	AppId      string
}

func main() {
	app := fiber.New()

	app.Get("/api/sse", adaptor.HTTPHandler(handler(statusHandler)))

	fmt.Println("Server listening on port 3000")
	app.Listen(":3000")
}

func handler(f http.HandlerFunc) http.Handler {
	return http.HandlerFunc(f)
}

func consumeOrderStatus(client *Client) {
	for {
		orderStatus := &OrderStatus{
			AppId:      "1",
			AppTransId: uuid.NewV4().String(),
		}

		client.events <- orderStatus

		time.Sleep(1 * time.Second)
	}
}

func statusHandler(w http.ResponseWriter, r *http.Request) {
	client := &Client{name: r.RemoteAddr, events: make(chan *OrderStatus, 1)}

	go consumeOrderStatus(client)

	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	timeout := time.After(1 * time.Minute)
	select {
	case ev := <-client.events:
		var buf bytes.Buffer
		enc := json.NewEncoder(&buf)
		enc.Encode(ev)
		fmt.Fprintf(w, "data: %v\n\n", buf.String())
		fmt.Printf("data: %v\n", buf.String())
	case <-timeout:
		fmt.Fprintf(w, ": nothing to sent\n\n")
	}

	if f, ok := w.(http.Flusher); ok {
		f.Flush()
	}
}
