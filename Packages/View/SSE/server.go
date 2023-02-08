package sse

import (
	"io"
	"log"

	"github.com/gin-gonic/gin"
)

// Broker It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier <-chan struct{}

	// New client connections
	newClients chan chan struct{}

	// Closed client connections
	closingClients chan chan struct{}

	// Total client connections
	clients map[chan struct{}]struct{}
}

// NewServer Initialize event and Start procnteessing requests
func NewServer() (b *Broker) {

	b = &Broker{
		newClients:     make(chan chan struct{}),
		closingClients: make(chan chan struct{}),
		clients:        make(map[chan struct{}]struct{}),
	}

	go b.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (b *Broker) listen() {
	for {
		select {
		case client := <-b.newClients:
			// Add new available client
			b.clients[client] = struct{}{}
			log.Printf("SSE Client added. %d registered clients", len(b.clients))

		case client := <-b.closingClients:
			// Remove closed client
			delete(b.clients, client)
			log.Printf("SSE Removed client. %d registered clients", len(b.clients))

		case eventMsg := <-b.Notifier:
			// Broadcast message to client
			for clientMessageChan := range b.clients {
				clientMessageChan <- eventMsg
			}
		}
	}
}

func (b *Broker) serveHTTP(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	// Initialize client channel
	messageChan := make(chan struct{})

	// Send new connection to event server
	b.newClients <- messageChan

	defer func() {
		b.closingClients <- messageChan
	}()

	notify := c.Request.Context().Done()
	go func() {
		<-notify
		b.closingClients <- messageChan
	}()

	defer func() {
		// Send closed connection to event server
		b.closingClients <- messageChan
	}()

	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		if msg, ok := <-messageChan; ok {
			c.SSEvent("update", msg)
			return true
		}
		return false
	})
}
