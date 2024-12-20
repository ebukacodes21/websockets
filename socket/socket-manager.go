package socket

import (
	"log"
	"sync"

	"github.com/gorilla/websocket"
)

type Client struct {
	Conn         *websocket.Conn
	Subscription string
}

type SocketManager struct {
	clients map[*websocket.Conn]*Client
	lock    sync.Mutex
}

func NewSocketManager() *SocketManager {
	return &SocketManager{clients: make(map[*websocket.Conn]*Client)}
}

func (m *SocketManager) AddClient(client *websocket.Conn, subscription string) {
	m.lock.Lock()
	defer m.lock.Unlock()

	m.clients[client] = &Client{
		Conn:         client,
		Subscription: subscription,
	}
	log.Printf("client added with address %s ", client.RemoteAddr())
}

func (m *SocketManager) RemoveClient(client *websocket.Conn) {
	m.lock.Lock()
	defer m.lock.Unlock()

	delete(m.clients, client)
	log.Printf("client with address %s removed", client.RemoteAddr())
}

func (m *SocketManager) BroadcastMessage(message interface{}, subscription string) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, client := range m.clients {
		if client.Subscription == subscription {
			if err := client.Conn.WriteJSON(message); err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Conn.Close()
				delete(m.clients, client.Conn)
			}
		}
	}

	return nil
}
