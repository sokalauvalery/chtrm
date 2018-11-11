package server

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"chat-room/internal/domain"
	"chat-room/internal/dto"

	"github.com/pkg/errors"
)

type ServerRestAPIController struct {
	chatServer domain.ServerManager
}

func (c ServerRestAPIController) AddClient(w http.ResponseWriter, r *http.Request) {
	client, err := readBodyClient(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	id, err := c.chatServer.RegisterClient(client.Name, client)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	loginResp := dto.LoginResponse{UUID: id}
	body, err := json.Marshal(loginResp)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Write(body)
}

// func (c ServerRestAPIController) DeleteClient(w http.ResponseWriter, r *http.Request) {
// 	client, err := readBodyClient(r.Body)
// 	if err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 	}

// 	if err = c.chatServer.DeleteClient(client.id); err != nil {
// 		http.Error(w, err.Error(), http.StatusInternalServerError)
// 		return
// 	}
// 	w.WriteHeader(http.StatusNoContent)
// }

func (c ServerRestAPIController) Broadcast(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		mgs := fmt.Sprintf("failed to read message request body %v", err)
		http.Error(w, mgs, http.StatusInternalServerError)
		return
	}

	messageBody := dto.MessageRequest{}

	if err = json.Unmarshal(b, &messageBody); err != nil {
		mgs := fmt.Sprintf("failed to unmarshal message request body %v", err)
		http.Error(w, mgs, http.StatusInternalServerError)
	}

	if err = c.chatServer.Broadcast(messageBody.UUID, messageBody.Message); err != nil {
		mgs := fmt.Sprintf("failed to broadcast message body %v", err)
		log.Printf("Broadcast error %v", mgs)
		http.Error(w, mgs, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func readBodyClient(body io.ReadCloser) (*Client, error) {
	b, err := ioutil.ReadAll(body)
	defer body.Close()
	if err != nil {
		return nil, errors.Wrap(err, "failed to read request body")
	}

	var client dto.LoginRequest
	if err = json.Unmarshal(b, &client); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal client from request body")
	}

	// return &Client{Name: client.Name, URL: client.URL}, nil
	return &Client{Name: client.Name,
		BaseURL:          client.BaseURL,
		MessageRoute:     client.MessageRoute,
		HealthcheckRoute: client.HealthcheckRoute}, nil
}
