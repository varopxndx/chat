package handler

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

// Processor methods
type Processor interface {
	ProcessMessage(message string) ([]byte, error)
}

// Broker methods
type Broker interface {
	GetQueueName() string
	Consume(key string) (<-chan amqp.Delivery, error)
}

// Socket struct
type Socket struct {
	clients    map[*Client]bool
	register   chan *Client
	unregister chan *Client
	broadcast  chan []byte
	processor  Processor
	rabbit     Broker
}

// NewWebsocketServer creates a new websocket server
func NewWebsocketServer(p Processor, r Broker) *Socket {
	return &Socket{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan []byte),
		processor:  p,
		rabbit:     r,
	}
}

func (s *Socket) startMQConsumer() {
	msgs, err := s.rabbit.Consume(s.rabbit.GetQueueName())
	if err != nil {
		panic(fmt.Sprintf("error to starting consumer: %v\n", err))
	}
	loop := make(chan bool)
	go func() {
		for msg := range msgs {
			for client := range s.clients {
				client.send <- msg.Body
			}
		}
	}()
	<-loop
}

// Run starts the consumer
func (s *Socket) Run() {
	go s.startMQConsumer()
	for {
		select {
		case client := <-s.register:
			s.registerClient(client)

		case client := <-s.unregister:
			s.unregisterClient(client)

		case message := <-s.broadcast:
			s.broadcastToClients(message)
		}
	}
}

func (s *Socket) registerClient(client *Client) {
	s.clients[client] = true
}

func (s *Socket) unregisterClient(client *Client) {
	delete(s.clients, client)
}

func (s *Socket) broadcastToClients(message []byte) {
	m, err := s.processor.ProcessMessage(string(message))
	if err != nil {
		log.Default().Printf("processing message error: %s", err.Error())
	}
	if m != nil {
		for client := range s.clients {
			client.send <- m
		}
	}
}
