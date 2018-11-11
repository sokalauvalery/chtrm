package domain

import (
	"fmt"
	"time"

	uuid "github.com/satori/go.uuid"
)

// TODO: parametrize
const defaultServerName = "server"

var defaultServerUUID = "00000000-0000-0000-0000-000000000000" // uuid.UUID{} //

// ServerManager interface with required server methods
type ServerManager interface {
	RegisterClient(string, ClientManager) (string, error)
	DeleteClient(string) error
	Broadcast(string, string) error
	ClientMonitor() error
}

type server struct {
	healthCheckInterval time.Duration
	clients             map[*ChatClient]ClientManager
	name                string
}

// NewServer builds new server instance
func NewServer() ServerManager {
	srv := server{
		name:                defaultServerName,
		clients:             make(map[*ChatClient]ClientManager),
		healthCheckInterval: time.Second * 1,
	}

	client := ChatClient{
		ID:   defaultServerUUID,
		Name: defaultServerName,
	}

	srv.clients[&client] = nil

	return &srv

}

// RegisterClient register client on server
func (srv *server) RegisterClient(name string, cm ClientManager) (string, error) {
	id, err := uuid.NewV4()
	if err != nil {
		return "", fmt.Errorf("failed to generate new uuid")
	}

	client := ChatClient{
		Name: name,
		ID:   id.String(),
	}

	srv.clients[&client] = cm
	if err = srv.Broadcast(defaultServerUUID, fmt.Sprintf("user %v online", client.Name)); err != nil {
		return "", fmt.Errorf("failed to register new user %v", err)
	}
	return id.String(), nil
}

// DeleteClient delete client from server
func (srv *server) DeleteClient(uuid string) error {
	client, err := srv.getClient(uuid)
	if err != nil {
		return fmt.Errorf("failed to delete client %v %v", uuid, err)
	}
	if _, ok := srv.clients[client]; ok {
		delete(srv.clients, client)
	}
	for c := range srv.clients {
		if c.Name == client.Name {
			// if there is any client with same name do not broadcast user offline message
			return nil
		}
	}
	if err = srv.Broadcast(defaultServerUUID, fmt.Sprintf("user %v offline", client.Name)); err != nil {
		return fmt.Errorf("failed to unregister user %v", err)
	}
	return nil
}

// Broadcast send message to all clients
func (srv *server) Broadcast(uuid, message string) error {
	sender, err := srv.getClient(uuid)
	if err != nil {
		return fmt.Errorf("failed to broadcast message: %v", err)
	}
	for client, manager := range srv.clients {
		if client.ID == uuid || client.ID == defaultServerUUID {
			continue
		}
		if err := manager.Send(fmt.Sprintf("%s: %s", sender.Name, message)); err != nil {
			return err
		}
	}
	return nil
}

// ClientMonitor check client state in infinite loop
func (srv *server) ClientMonitor() error {
	for {
		for client, manager := range srv.clients {
			if client.ID != defaultServerUUID && !manager.HealthCheck() {
				if err := srv.DeleteClient(client.ID); err != nil {
					return fmt.Errorf("failed to delete client %v", err)
				}
			}
		}
		time.Sleep(srv.healthCheckInterval)
	}
}

func (srv *server) getClient(id string) (*ChatClient, error) {
	for client := range srv.clients {
		if client.ID == id {
			return client, nil
		}
	}
	return nil, fmt.Errorf("client with uuid %v not foud", id)
}
