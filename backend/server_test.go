package main

import (
	"github.com/gorilla/websocket"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTurnOrder(t *testing.T) {
	room := NewChatRoom(nil)

	player1 := &websocket.Conn{}
	player2 := &websocket.Conn{}

	room.turnOrder = append(room.turnOrder, player1, player2)

	if room.turnOrder[0] != player1 {
		t.Errorf("Ошибка: первым должен быть player1")
	}

	room.turnOrder = append(room.turnOrder[1:], player1)

	if room.turnOrder[0] != player2 {
		t.Errorf("Ошибка: первым теперь должен быть player2")
	}
}

func TestCreateRoom(t *testing.T) {
	server := NewServer(nil)

	req := httptest.NewRequest("GET", "/create", nil)
	w := httptest.NewRecorder()

	server.CreateRoom(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusOK {
		t.Errorf("Ошибка: ожидался статус 200, получен %d", resp.StatusCode)
	}
}

func TestJoinRoom_Fail(t *testing.T) {
	server := NewServer(nil)

	req := httptest.NewRequest("GET", "/ws/неизвестная-комната", nil)
	w := httptest.NewRecorder()

	server.JoinRoom(w, req)

	resp := w.Result()
	if resp.StatusCode != http.StatusNotFound {
		t.Errorf("Ошибка: ожидался статус 404, получен %d", resp.StatusCode)
	}
}
