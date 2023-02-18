package sse

import (
	"io"
	"log"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"chicha/Packages/race"
)

// Client хранит raceID гонки и канал для отправки уведомлений клиенту
type Client struct {
	raceID race.ID
	notify chan struct{}
}

// Broker It keeps a list of clients those are currently attached
// and broadcasting events to those clients.
type Broker struct {
	// Events are pushed to this channel by the main events-gathering routine
	Notifier <-chan race.ID

	// New client connections
	newClients chan Client

	// Closed client connections
	closingClients chan Client

	// Total client connections
	clients map[race.ID]map[Client]struct{} // все текущие гонки и клиенты подключенные к ним
}

// NewServer Initialize event and Start proceessing requests
func NewServer() (b *Broker) {
	b = &Broker{
		newClients:     make(chan Client),
		closingClients: make(chan Client),
		clients:        make(map[race.ID]map[Client]struct{}),
	}

	go b.listen()

	return
}

// It Listens all incoming requests from clients.
// Handles addition and removal of clients and broadcast messages to clients.
func (b *Broker) listen() {
	for {
		select {
		case newClient := <-b.newClients:
			// если это первый клиент для гонки, то мапа будет nil, инициализируем её
			if b.clients[newClient.raceID] == nil {
				b.clients[newClient.raceID] = make(map[Client]struct{})
			}

			// добавляем клиента к гонке
			b.clients[newClient.raceID][newClient] = struct{}{}
			log.Printf("SSE Client added. %d registered clients for race %v", len(b.clients[newClient.raceID]), newClient.raceID)

		case client := <-b.closingClients:
			if b.clients[client.raceID] == nil {
				continue
			}
			log.Printf("SSE Removed client. %d registered clients for race %v", len(b.clients[client.raceID]), client.raceID)

			// Remove closed client
			delete(b.clients[client.raceID], client)
			log.Printf("SSE Removed client. %d registered clients for race %v", len(b.clients[client.raceID]), client.raceID)

		case raceID := <-b.Notifier:
			// ищем подписаных клиентов под raceID гонки, и отправляем только тех, кто подписан на raceID гонки из уведомления
			if b.clients[raceID] == nil {
				continue
			}
			for client := range b.clients[raceID] {
				client.notify <- struct{}{}
			}
		}
	}
}

func (b *Broker) serveHTTP(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "text/event-stream")
	c.Writer.Header().Set("Cache-Control", "no-cache")
	c.Writer.Header().Set("Connection", "keep-alive")
	c.Writer.Header().Set("Transfer-Encoding", "chunked")

	raceID, _ := strconv.ParseUint(c.Params.ByName("raceID"), 10, 64)
	log.Println("RACE ID", raceID)

	// Initialize client
	newClient := Client{
		raceID: race.ID(raceID),
		notify: make(chan struct{}, 1), // на всякий случай буфер чтобы не блочилось
	}

	// Send new connection to event server
	b.newClients <- newClient

	defer func() {
		// удаляем клиентов после того как они уйдут со страницы и закроют соединение
		// Send closed connection to event server
		b.closingClients <- newClient
	}()

	timer := time.NewTimer(time.Second * 60) // открываем таймер для проверки подключения клиентов
	defer timer.Stop()                       // закрываем таймер чтобы не было утечек
	c.Stream(func(w io.Writer) bool {
		// Stream message to client from message channel
		select {
		case msg := <-newClient.notify:
			c.SSEvent("update", msg)
			timer.Reset(time.Second * 60) // обновляем таймер если получилось отправить сообщение
			return true
		case <-timer.C:
			return false
		}
	})
	log.Println("CONN CLOSED")
}
