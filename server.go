package main

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"net/http"
	"sync"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Client struct {
	conn      *websocket.Conn
	connected bool
}

type ChatRoom struct {
	messages    map[string]int
	lastMessage string
	clients     map[*websocket.Conn]*Client
	broadcast   chan []byte
	mutex       sync.Mutex
	turnOrder   []*websocket.Conn
	wordChecker *WordCheckerClient
}

func NewChatRoom(wordChecker *WordCheckerClient) *ChatRoom {
	room := &ChatRoom{
		clients:     make(map[*websocket.Conn]*Client),
		broadcast:   make(chan []byte),
		messages:    make(map[string]int),
		wordChecker: wordChecker,
	}
	go room.Run()
	return room
}

func (cr *ChatRoom) Run() {
	for msg := range cr.broadcast {
		cr.mutex.Lock()
		for _, client := range cr.clients {
			if err := client.conn.WriteMessage(websocket.TextMessage, msg); err != nil {
				client.conn.Close()
				delete(cr.clients, client.conn)
			}
		}
		cr.mutex.Unlock()
	}
}

type Server struct {
	rooms       map[string]*ChatRoom
	mutex       sync.Mutex
	wordChecker *WordCheckerClient
}

func NewServer(wordChecker *WordCheckerClient) *Server {
	return &Server{rooms: make(map[string]*ChatRoom), wordChecker: wordChecker}
}

func (s *Server) CreateRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	roomID := uuid.New().String()

	s.mutex.Lock()
	s.rooms[roomID] = NewChatRoom(s.wordChecker)
	s.mutex.Unlock()

	roomURL := fmt.Sprintf("ws://%s/ws/%s", r.Host, roomID)
	w.Write([]byte(roomURL))
}

func (s *Server) JoinRoom(w http.ResponseWriter, r *http.Request) {
	roomID := r.URL.Path[len("/ws/"):]

	s.mutex.Lock()
	room, exists := s.rooms[roomID]
	s.mutex.Unlock()

	if !exists {
		http.Error(w, "Room not found", http.StatusNotFound)
		return
	}

	room.HandleConnection(w, r)
}
