package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"chat-room/internal/dto"

	"github.com/pkg/errors"
	uuid "github.com/satori/go.uuid"
)

const msgContentType = "application/json"

// Client contains data related to chat client
type Client struct {
	id   uuid.UUID
	Name string

	BaseURL          string
	MessageRoute     string
	HealthcheckRoute string
}

// Send message to client
func (c Client) Send(msg string) error {
	client := http.Client{Timeout: 5 * time.Second}
	dtoMessage := dto.MessageResponse{
		Time: time.Now(),
		Text: msg,
	}

	bts, err := json.Marshal(&dtoMessage)
	if err != nil {
		return err
	}

	messageURL := fmt.Sprintf("%v/%v", c.BaseURL, c.MessageRoute)
	log.Printf("Send message to %s", messageURL)
	resp, err := client.Post(messageURL, msgContentType, bytes.NewReader(bts))
	if err != nil {
		log.Printf("ERROR %s", err)
		// return errors.Wrapf(err, "cannot send message to client %s-%s", c.Name, c.ID)
	}

	if resp.StatusCode != http.StatusNoContent {
		return fmt.Errorf("client returns %v status code", resp.StatusCode)
	}

	return nil
}

// HealthCheck check is client online
func (c Client) HealthCheck() (state bool) {
	client := http.Client{Timeout: 1 * time.Second}

	log.Printf("%v", fmt.Sprintf("%v/%v", c.BaseURL, c.HealthcheckRoute))
	resp, err := client.Get(fmt.Sprintf("%v/%v", c.BaseURL, c.HealthcheckRoute))
	if err != nil {
		log.Print(errors.Wrapf(err, "healthcheck failed"))
		return
	}

	if resp.StatusCode != http.StatusNoContent {
		log.Print(fmt.Errorf("client returns %v status code", resp.StatusCode))
		return
	}

	return true
}
