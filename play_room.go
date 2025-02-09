package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
)

func (cr *ChatRoom) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, _ := upgrader.Upgrade(w, r, nil)

	if !cr.addClient(conn) {
		conn.WriteMessage(websocket.TextMessage, []byte("Комната уже заполнена!"))
		conn.Close()
		return
	}

	defer cr.removeClient(conn)

	cr.notifyClients()
	cr.listenForMessages(conn)
}

func (cr *ChatRoom) addClient(conn *websocket.Conn) bool {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	if len(cr.clients) >= 2 {
		return false
	}

	client := &Client{
		conn:      conn,
		connected: true,
	}
	cr.clients[conn] = client
	cr.turnOrder = append(cr.turnOrder, conn)
	return true
}

func (cr *ChatRoom) removeClient(conn *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	delete(cr.clients, conn)
	for i, c := range cr.turnOrder {
		if c == conn {
			cr.turnOrder = append(cr.turnOrder[:i], cr.turnOrder[i+1:]...)
			break
		}
	}
	conn.Close()
}

func (cr *ChatRoom) notifyClients() {
	msg := fmt.Sprintf("{\"connected\": %d}", len(cr.clients))
	cr.broadcast <- []byte(msg)
}

func (cr *ChatRoom) listenForMessages(conn *websocket.Conn) {
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println(err)
			break
		}
		cr.processMessage(conn, msg)
	}
}

func (cr *ChatRoom) processMessage(conn *websocket.Conn, msg []byte) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	var response GameResponse
	if cr.turnOrder[0] == conn {
		if err := json.Unmarshal(msg, &response); err != nil {
			fmt.Println("Ошибка при разборе JSON:", err)
			return
		}

		switch response.Type {
		case "timeout":
			cr.handleTimeout(conn)
		case "turn":
			cr.handleTurn(conn)
		default:
			cr.handleWordMove(conn, response.Text)
		}
	}
}

func (cr *ChatRoom) handleTimeout(loser *websocket.Conn) {
	var winner *websocket.Conn
	for _, c := range cr.turnOrder {
		if c != loser {
			winner = c
			break
		}
	}
	if loser != nil {
		loser.WriteMessage(websocket.TextMessage, []byte(`{"type":"result","text":"Вы проиграли"}`))
		loser.Close()
		delete(cr.clients, loser)
	}
	if winner != nil {
		winner.WriteMessage(websocket.TextMessage, []byte(`{"type":"result","text":"Вы победили"}`))
		winner.Close()
		delete(cr.clients, winner)
	}
	cr.turnOrder = nil
}

func (cr *ChatRoom) handleTurn(secondTurn *websocket.Conn) {
	var firstTurn *websocket.Conn
	for _, c := range cr.turnOrder {
		if c != secondTurn {
			firstTurn = c
			break
		}
	}

	if secondTurn != nil {
		secondTurn.WriteMessage(websocket.TextMessage, []byte(`{"type":"turn","text":"Ваш Ход!"}`))
	}

	if firstTurn != nil {
		firstTurn.WriteMessage(websocket.TextMessage, []byte(`{"type":"turn","text":"Ход Соперника!"}`))
	}
}

func (cr *ChatRoom) handleWordMove(conn *websocket.Conn, receivedWord string) {
	isMoveCorrect, err := cr.wordChecker.IsMoveCorrect(cr.lastMessage, receivedWord, cr.messages)
	if err != nil {
		fmt.Println(err)
		msg := fmt.Sprintf("{\"type\": \"error\", \"text\": \"%s\"}", err.Error())
		cr.turnOrder[0].WriteMessage(websocket.TextMessage, []byte(msg))
		return
	}

	if isMoveCorrect {
		cr.messages[receivedWord] = 0
		cr.broadcast <- []byte(receivedWord)

		msg := fmt.Sprintf("{\"type\": \"turn\", \"text\": \"%s\"}", "Ход Соперника!")
		cr.turnOrder[0].WriteMessage(websocket.TextMessage, []byte(msg))

		cr.lastMessage = receivedWord
		cr.turnOrder = append(cr.turnOrder[1:], conn)

		msg = fmt.Sprintf("{\"type\": \"turn\", \"text\": \"%s\"}", "Ваш Ход!")
		cr.turnOrder[0].WriteMessage(websocket.TextMessage, []byte(msg))
	}
}
