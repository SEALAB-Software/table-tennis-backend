// Copyright 2013 The Gorilla WebSocket Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package main

import (
	"fmt"
)


// Hub maintains the set of active clients and broadcasts messages to the
// clients.
type Hub struct {
	// Registered clients.
	clients map[*Client]bool

	// Inbound messages from the clients.
	broadcast chan []byte

	// Register requests from the clients.
	register chan *Client

	// Unregister requests from clients.
	unregister chan *Client

	// Increment score_1
	increment_score_1 chan int

	score_1 int

	// Increment score_2
	increment_score_2 chan int

	score_2 int
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
		increment_score_1: 	make(chan int),
		increment_score_2: 	make(chan int),
		score_1: 0,
		score_2: 0,
	}
}

func publish_message(h *Hub, message string) int {
	b := []byte(message)
	for client := range h.clients {
		select {
		case client.send <- b:
		default:
			fmt.Println("Failed to send client message!")
			close(client.send)
			delete(h.clients, client)
		}
	}
	return 0
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case increment := <-h.increment_score_1:
			h.score_1 += increment
			fmt.Println(h.score_1)
			jsonString := fmt.Sprintf("{\"score_1\": %d, \"score_2\": %d}", h.score_1, h.score_2)
			publish_message(h, jsonString)
		case increment := <-h.increment_score_2:
			h.score_2 += increment
			fmt.Println(h.score_2)
			jsonString := fmt.Sprintf("{\"score_1\": %d, \"score_2\": %d}", h.score_1, h.score_2)
			publish_message(h, jsonString)
		}
	}
}
	// case message := <-h.broadcast:
	// 	for client := range h.clients {
	// 		select {
	// 		case client.send <- message:
	// 		default:
	// 			close(client.send)
	// 			delete(h.clients, client)
	// 		}
	// 	}
	// }
