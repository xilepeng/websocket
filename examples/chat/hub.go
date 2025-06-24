package main

import "log"

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte  // 广播消息
	register   chan *Client // 注册
	unregister chan *Client // 注销
}

func newHub() *Hub {
	return &Hub{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}
}

// 服务端：监听并广播消息
func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			log.Println("Register client:", client)
			h.clients[client] = true
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			log.Println("Broadcast message:", string(message))
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(h.clients, client)
				}
			}
		}
	}
}
