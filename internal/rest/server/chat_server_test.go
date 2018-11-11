package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"chat-room/internal/domain"

	"github.com/stretchr/testify/require"
)

type mockServerManager struct {
	clientsToRegister   []domain.ChatClient
	clientsToDelete     []domain.ChatClient
	messagesToBroadCast []string
}

func (srv *mockServerManager) RegisterClient(client *domain.ChatClient) error {
	srv.clientsToRegister = append(srv.clientsToRegister, *client)
	return nil
}
func (srv *mockServerManager) DeleteClient(id string) error {
	toDelete := domain.ChatClient{
		ID: id,
	}
	srv.clientsToDelete = append(srv.clientsToDelete, toDelete)
	return nil
}
func (srv *mockServerManager) Broadcast(uuid, message string) error {
	srv.messagesToBroadCast = append(srv.messagesToBroadCast, message)
	return nil
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected response code %d. Got %d\n", expected, actual)
	}
}

func TestRestGateWay(t *testing.T) {

	srv := &mockServerManager{}
	webSerAddr := "127.0.0.1:1234"
	testMessage := "test message"
	testChatClient := domain.Client{Name: "1"}
	srvGate, _ := NewChatWebServer(webSerAddr, srv)

	go func() {
		err := srvGate.Run()
		require.Nil(t, err, "Server run terminated")
	}()

	client := http.Client{}
	defaultContentType := "application/json"

	clientBody, _ := json.Marshal(testChatClient)
	resp, err := client.Post(fmt.Sprintf("http://%s/%s", webSerAddr, "login"), defaultContentType, bytes.NewReader(clientBody))
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Equal(t, 1, len(srv.clientsToRegister))
	require.Equal(t, testChatClient.Name, srv.clientsToRegister[0].Name)

	resp, err = client.Post(fmt.Sprintf("http://%s/%s", webSerAddr, "logout"), defaultContentType, bytes.NewReader(clientBody))
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Equal(t, 1, len(srv.clientsToDelete), fmt.Sprintf("%v", srv.clientsToDelete))
	require.Equal(t, testChatClient.Name, srv.clientsToDelete[0].Name)

	resp, err = client.Post(fmt.Sprintf("http://%s/%s", webSerAddr, "message"), defaultContentType, bytes.NewReader([]byte(testMessage)))
	require.Nil(t, err)
	require.Equal(t, http.StatusNoContent, resp.StatusCode)
	require.Equal(t, 1, len(srv.messagesToBroadCast), fmt.Sprintf("%v", srv.messagesToBroadCast))
	require.Equal(t, testMessage, srv.messagesToBroadCast[0])

}
