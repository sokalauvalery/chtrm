package client

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"chat-room/internal/dto"
)

// AgentRestAPIController chat client rest api controller
type AgentRestAPIController struct {
	chatClient ChatClient
}

// Message incoming message receive handler
func (c AgentRestAPIController) Message(w http.ResponseWriter, r *http.Request) {
	b, err := ioutil.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		mgs := fmt.Sprintf("failed to read message request body %v", err)
		http.Error(w, mgs, http.StatusInternalServerError)
		return
	}

	messageBody := dto.MessageResponse{}

	if err = json.Unmarshal(b, &messageBody); err != nil {
		mgs := fmt.Sprintf("failed to unmarshal message request body %v", err)
		http.Error(w, mgs, http.StatusInternalServerError)
	}

	if err = c.chatClient.Receive(messageBody.Text); err != nil {
		mgs := fmt.Sprintf("failed to recieve message body %v", err)
		log.Printf("Receive error %v", mgs)
		http.Error(w, mgs, http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Healthcheck updates last connection state
func (c AgentRestAPIController) Healthcheck(w http.ResponseWriter, r *http.Request) {
	c.chatClient.UpdateLastConnectionTime()
	w.WriteHeader(http.StatusNoContent)
}
