package utils

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Room struct {
	clients   map[*websocket.Conn]bool
	broadcast chan []byte
	lock      sync.Mutex
}

type SocketServer struct {
	rooms map[string]*Room
	lock  sync.Mutex
}

type SocketMessage struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

func NewSocketServer() *SocketServer {
	return &SocketServer{
		rooms: make(map[string]*Room),
	}
}

var socketServer = NewSocketServer()

func (s *SocketServer) JoinRoom(roomName string, conn *websocket.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()

	room, exists := s.rooms[roomName]
	if !exists {
		room = &Room{
			clients:   make(map[*websocket.Conn]bool),
			broadcast: make(chan []byte),
		}
		s.rooms[roomName] = room
		go s.handleMessages(room) // Start listening for messages in this room
	}

	room.lock.Lock()
	room.clients[conn] = true
	room.lock.Unlock()
}

func (s *SocketServer) LeaveRoom(roomName string, conn *websocket.Conn) {
	s.lock.Lock()
	defer s.lock.Unlock()

	if room, exists := s.rooms[roomName]; exists {
		room.lock.Lock()
		delete(room.clients, conn)
		room.lock.Unlock()

		if len(room.clients) == 0 {
			delete(s.rooms, roomName)
		}
	}
}

func (s *SocketServer) handleMessages(room *Room) {
	for {
		msg := <-room.broadcast
		room.lock.Lock()
		for client := range room.clients {
			err := client.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				log.Printf("Error sending message to client: %v", err)
				client.Close()
				delete(room.clients, client)
			}
		}
		room.lock.Unlock()
	}
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

func (s *SocketServer) HandleWebSocket(c *gin.Context) {
	roomName := c.Param("room")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("Failed to upgrade connection:", err)
		return
	}
	defer conn.Close()

	s.JoinRoom(roomName, conn)
	defer s.LeaveRoom(roomName, conn)

	for {
		// No handling of incoming messages
		_, _, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			break
		}
	}
}

func SendMessageToRoom(roomName string, data SocketMessage) error {
	socketServer.lock.Lock()
	defer socketServer.lock.Unlock()

	room, exists := socketServer.rooms[roomName]
	if !exists {
		return fmt.Errorf("room %s does not exist", roomName)
	}

	message, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("error marshaling data to JSON: %v", err)
	}

	room.broadcast <- message
	return nil
}

// room format: restaurantId_tableId
func InitWebSocketServer(r *gin.Engine) {
	r.GET("/:room", func(c *gin.Context) {
		socketServer.HandleWebSocket(c)
	})
}
