package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

type client struct {
	conn *websocket.Conn
	mu   sync.Mutex
}

var upgrader = websocket.Upgrader{
	Subprotocols: []string{"websocket"},
	CheckOrigin: func(r *http.Request) bool {
		if r.Header.Get("Origin") == "http://localhost:8080" || r.Header.Get("Origin") == "http://localhost:8081" {
			return true
		}
		return false
	},
}
var clients = make(map[string]*client)
var clientsMu sync.Mutex

func main() {
	// Connect to Cassandra DB
	cluster := gocql.NewCluster("127.0.0.1")
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatal(err)
	}
	defer session.Close()

	// Set up WebSocket server
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Fatal(err)
		}

		// Store connection in clients map
		var userID string
		fmt.Sscanf(r.URL.Query().Get("userID"), "%s", &userID)
		clientsMu.Lock()
		clients[userID] = &client{conn: conn}
		clientsMu.Unlock()

		// Handle incoming WebSocket messages
		for {
			messageType, p, err := conn.ReadMessage()
			if err != nil {
				log.Println(err)
				return
			}
			if messageType != websocket.TextMessage {
				continue
			}

			// Parse WebSocket message
			var notificationText string
			fmt.Sscanf(string(p), "%s", &notificationText)

			// Query Cassandra DB for recipient's name
			var recipientName string
			query := "SELECT name FROM notifications.users WHERE id = ?"
			err = session.Query(query, userID).Consistency(gocql.One).Scan(&recipientName)
			if err != nil {
				log.Println(err)
				return
			}

			// Construct notification message
			notificationMessage := fmt.Sprintf("Dear %s,\n\nYou have a new notification:\n%s", recipientName, notificationText)

			// Send notification message over WebSocket connection
			err = sendMessageToClient(userID, notificationMessage)
			if err != nil {
				log.Println(err)
				return
			}

		}
	})

	// Route for sending messages to clients
	http.HandleFunc("/send", func(w http.ResponseWriter, r *http.Request) {
		userID := r.URL.Query().Get("userID")
		message := r.FormValue("message")

		err := sendMessageToClient(userID, message)
		if err != nil {
			log.Println(err)
			http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
	})

	// Start server
	log.Print("Starting server at ws://localhost:8081/ws and http://localhost:8081/send")
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func sendMessageToClient(userID string, message string) error {
	// Get the client's WebSocket connection from the clients map
	clientsMu.Lock()
	client := clients[userID]
	clientsMu.Unlock()

	if client == nil {
		return fmt.Errorf("client with userID %s not found", userID)
	}

	// Lock the client's mutex to prevent concurrent writes to the connection
	client.mu.Lock()
	defer client.mu.Unlock()

	// Send the message over the WebSocket connection
	err := client.conn.WriteMessage(websocket.TextMessage, []byte(message))
	if err != nil {
		return err
	}

	return nil
}
