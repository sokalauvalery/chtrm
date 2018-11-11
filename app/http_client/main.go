package main

import (
	"log"
	"time"

	"chat-room/internal/rest/client"
	"chat-room/internal/ui/console"

	"github.com/pkg/errors"
)

var serverStateCheckInterval = 1 * time.Second

// TODO: move logic to libraries
func main() {
	serverURL := "http://127.0.0.1:4444"

	term := console.New()
	if err := term.WriteToUser("Enter your nickname: "); err != nil {
		log.Fatal(err)
	}
	username, err := term.ReadInput()
	if err != nil {
		log.Fatal(err)
	}

	inbox := make(chan string)

	go func() {
		if err = printIncomingMessages(term, inbox); err != nil {
			log.Fatal(err)
		}
	}()

	chatClient, err := client.NewChatClient(username, serverURL, inbox)
	if err != nil {
		log.Fatal(err)
	}

	go func() {
		log.Fatal(chatClient.StartServer())
	}()

	if err = chatClient.Login(); err != nil {
		log.Fatal(err)
	}

	go func() {
		serverConnectionMonitor(term, chatClient)
	}()

	for {
		msg, err := term.ReadInput()
		if err != nil {
			log.Fatal(err)
		}

		if err = chatClient.Send(msg); err != nil {
			log.Fatal(err)
		}
	}

}

func printIncomingMessages(term console.UI, inbox chan string) error {
	var message string
	for {
		message = <-inbox
		if err := term.WriteToUser(message); err != nil {
			return errors.Wrap(err, "failed to show message to user")
		}
	}
}

func serverConnectionMonitor(term console.UI, client client.ChatClient) {
	for {
		client.CheckConnection()
		time.Sleep(serverStateCheckInterval)
	}
}
