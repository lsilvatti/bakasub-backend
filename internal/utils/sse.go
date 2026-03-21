package utils

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type SSEEvent struct {
	Type    string         `json:"type"`
	Module  string         `json:"module"`
	Message string         `json:"message"`
	Data    map[string]any `json:"data,omitempty"`
}

type SSEBroker struct {
	clients        map[chan SSEEvent]bool
	newClients     chan chan SSEEvent
	defunctClients chan chan SSEEvent
	messages       chan SSEEvent
}

var Broker *SSEBroker

func InitSSEBroker() {
	Broker = &SSEBroker{
		clients:        make(map[chan SSEEvent]bool),
		newClients:     make(chan chan SSEEvent),
		defunctClients: make(chan chan SSEEvent),
		messages:       make(chan SSEEvent),
	}
	go Broker.start()
}

func (b *SSEBroker) start() {
	for {
		select {
		case s := <-b.newClients:
			b.clients[s] = true
		case s := <-b.defunctClients:
			delete(b.clients, s)
			close(s)
		case msg := <-b.messages:
			for s := range b.clients {
				select {
				case s <- msg:
				default:
					delete(b.clients, s)
					close(s)
				}
			}
		}
	}
}

func SendSSE(eventType, module, message string, data map[string]any) {
	if Broker != nil {
		Broker.messages <- SSEEvent{
			Type:    eventType,
			Module:  module,
			Message: message,
			Data:    data,
		}
	}
}

func (b *SSEBroker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported!", http.StatusInternalServerError)
		return
	}

	messageChan := make(chan SSEEvent, 10)
	b.newClients <- messageChan

	notify := r.Context().Done()
	go func() {
		<-notify
		b.defunctClients <- messageChan
	}()

	for {
		msg, open := <-messageChan
		if !open {
			break
		}

		jsonData, _ := json.Marshal(msg)
		fmt.Fprintf(w, "data: %s\n\n", jsonData)
		flusher.Flush()
	}
}
