package main

import (
	"log"
	"net/http"
	"os"

	"github.com/DylanCoon99/chatapp/trace"
	"github.com/gorilla/websocket"
)

type room struct {
	forward chan []byte      // holds all incoming messages
	join    chan *client     // holds all clients wishing to join the room
	leave   chan *client     // holds all clients wanting to leave the room
	clients map[*client]bool // holds all clients
	tracer  trace.Tracer
}

func (r *room) run() {
	r.tracer.Trace("Room is open!")
	for {
		select {
		case client := <-r.join:
			// a client wants to join the room
			r.clients[client] = true
			r.tracer.Trace("New Client Joined")
		case client := <-r.leave:
			// a client wants to leave the room
			delete(r.clients, client)
			close(client.send)
			r.tracer.Trace("Client Left")
		case msg := <-r.forward:
			// there is an incoming message that needs to be sent to all clients
			r.tracer.Trace("Message Received: ", string(msg))
			for client := range r.clients {
				client.send <- msg
				r.tracer.Trace(" -- sent to client")
			}
		}
	}
}

const (
	socketBufferSize  = 1024
	messageBufferSize = 256
)

var upgrader = &websocket.Upgrader{ReadBufferSize: socketBufferSize, WriteBufferSize: socketBufferSize}

func (r *room) ServeHTTP(w http.ResponseWriter, req *http.Request) {

	socket, err := upgrader.Upgrade(w, req, nil)

	if err != nil {
		log.Fatal("ServeHTTP: ", err)
		return
	}

	client := &client{
		socket: socket,
		send:   make(chan []byte, messageBufferSize),
		room:   r,
	}

	r.join <- client

	defer func() { r.leave <- client }()
	go client.write()
	client.read()

}

func newRoom() *room {
	return &room{
		forward: make(chan []byte),
		join:    make(chan *client),
		leave:   make(chan *client),
		clients: make(map[*client]bool),
		tracer:  trace.New(os.Stdout),
	}
}
