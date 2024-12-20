package servers

import (
	"context"
	"log"
	"net/http"
	"time"
	"websockets/socket"

	"github.com/gorilla/websocket"
	"golang.org/x/sync/errgroup"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func RunWebsocket(ctx context.Context, group *errgroup.Group, manager *socket.SocketManager) {
	mux := http.NewServeMux()

	mux.HandleFunc("/ws/update-order", func(w http.ResponseWriter, r *http.Request) {
		handleUpdateOrder(w, r, manager, "update-order")
	})

	server := &http.Server{
		Addr:    "0.0.0.0:8000",
		Handler: mux,
	}

	group.Go(func() error {
		log.Print("running....")
		err := server.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}

		return nil
	})

	group.Go(func() error {
		<-ctx.Done()
		log.Print("gracefully shutting down...")

		err := server.Shutdown(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		log.Print("goodbye")
		return nil
	})

	go broadcastOrder(manager)
}

func handleUpdateOrder(w http.ResponseWriter, r *http.Request, manager *socket.SocketManager, subcription string) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Fatal(err)
	}

	defer conn.Close()

	manager.AddClient(conn, subcription)

	// keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			manager.RemoveClient(conn)
			break
		}
	}
}

func broadcastOrder(manager *socket.SocketManager) {
	ticker := time.NewTicker(2 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			newOrderMessage := map[string]interface{}{
				"orderID":   12345,
				"status":    "created",
				"message":   "A new order has been created.",
				"createdAt": time.Now(),
			}

			err := manager.BroadcastMessage(newOrderMessage, "update-order")
			if err != nil {
				log.Printf("Error broadcasting new order: %v", err)
			}
		}
	}
}
