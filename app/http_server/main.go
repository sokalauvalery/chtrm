package main

import (
	"log"

	"chat-room/internal/domain"
	"chat-room/internal/rest/server"
)

func main() {
	srv := domain.NewServer()
	srvGate, _ := server.NewChatWebServer("127.0.0.1:4444", srv)
	log.Fatal(srvGate.Run())
}
