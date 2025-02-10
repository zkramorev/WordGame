package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	wordChecker, err := NewWordCheckerClient("wordgame-python-server:50051")
	if err != nil {
		log.Fatalf("Ошибка при подключении к gRPC-серверу: %v", err)
	}
	defer wordChecker.Close()

	server := NewServer(wordChecker)

	http.HandleFunc("/create", server.CreateRoom)
	http.HandleFunc("/ws/", server.JoinRoom)

	fmt.Println("Server started on :8080")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		fmt.Println("Error starting server:", err)
	}
}
