package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"chat-room/internal/dto"
	"chat-room/internal/rest/common"

	"github.com/pkg/errors"
)

const clientContentType = "application/json"

var serverHealthcheckTimeout = 5 * time.Second

// ChatClient chat client interface
type ChatClient interface {
	// TODO: none-gopher huge interface - use couple of small
	Login() error

	Send(string) error
	Receive(string) error
	StartServer() error

	String() string
	UpdateLastConnectionTime()
	CheckConnection()
}

type chatClient struct {
	username string
	uuid     string
	ws       *common.WebServer

	clientURL  string
	serverURL  string
	httpClient http.Client

	inbox chan string

	lastConnection time.Time
	online         bool
}

func (c chatClient) String() string {
	return c.username + " " + c.uuid
}

// Login establishes connection to server
func (c *chatClient) Login() error {
	loginURL := fmt.Sprintf("%s/login", c.serverURL)
	loginBody := dto.LoginRequest{
		Name:             c.username,
		BaseURL:          fmt.Sprintf("%s", c.clientURL),
		MessageRoute:     "message",
		HealthcheckRoute: "healthcheck",
	}
	body, err := json.Marshal(&loginBody)
	if err != nil {
		return errors.Wrap(err, "failed to marshal request body")
	}
	resp, err := c.httpClient.Post(loginURL, clientContentType, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "failed to post login to server")
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("connection to server failed %v", resp.StatusCode)
	}

	loginResp := dto.LoginResponse{}
	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read response body")
	}

	if err = json.Unmarshal(bodyBytes, &loginResp); err != nil {
		return errors.Wrap(err, "failed to unmarshal login response body")
	}
	c.uuid = loginResp.UUID
	c.online = true
	return nil
}

// Send sends message to server to broadcast for other users
func (c *chatClient) Send(msg string) error {
	if !c.online {
		c.inbox <- "server is offline, try later"
		return nil
	}
	sendMessageURL := fmt.Sprintf("%s/message", c.serverURL)
	msgBody := dto.MessageRequest{
		UUID:    c.uuid,
		Message: msg,
	}
	body, err := json.Marshal(&msgBody)
	if err != nil {
		return errors.Wrap(err, "failed to marshal message request body")
	}
	resp, err := c.httpClient.Post(sendMessageURL, clientContentType, bytes.NewReader(body))
	if err != nil {
		return errors.Wrap(err, "failed to post login to server")
	}

	if resp.StatusCode != http.StatusNoContent {
		defer resp.Body.Close()
		bodyBytes, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("failed to send message %v", resp.StatusCode)
		}
		return fmt.Errorf("failed to send message - code %v msg %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// StartServer starts chat client server
func (c *chatClient) StartServer() error {
	return c.ws.ListenAndServe()
}

// Receive shows incoming message
func (c *chatClient) Receive(msg string) error {
	c.inbox <- msg
	return nil
}

// UpdateLastConnectionTime
func (c *chatClient) UpdateLastConnectionTime() {
	c.lastConnection = time.Now()
}

func (c *chatClient) CheckConnection() {
	if c.lastConnection.IsZero() {
		return
	}
	if time.Since(c.lastConnection) > serverHealthcheckTimeout && c.online {
		c.inbox <- fmt.Sprintf("connection to the server lost. last connection at %v", c.lastConnection)
		c.online = false
		return
	}
	if time.Since(c.lastConnection) < serverHealthcheckTimeout && !c.online {
		c.inbox <- "server connection restored"
		c.online = true
		return
	}

}

// NewChatClient builds new chat client object
func NewChatClient(username, serverURL string, inbox chan string) (ChatClient, error) {
	httpClient := http.Client{Timeout: 1 * time.Second}

	listener, err := net.Listen("tcp", ":0")
	if err != nil {
		return nil, err
	}

	ws := common.NewWebServerWithListener(&listener)

	client := &chatClient{
		username:   username,
		httpClient: httpClient,
		serverURL:  serverURL,
		ws:         ws,
		clientURL:  fmt.Sprintf("http://%v", listener.Addr().String()),
		inbox:      inbox,
	}

	routes := GetChatClientRoutes(client)
	for _, route := range routes {
		client.ws.AddRoute(route.Method, route.URL, route.Handler)
	}

	return client, nil

}
