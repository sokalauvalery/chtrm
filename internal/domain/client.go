package domain

// ClientManager interface for interaction with server client
type ClientManager interface {
	// Name() string
	Send(string) error
	HealthCheck() bool
}

// ChatClient client indentification info
type ChatClient struct {
	ID   string
	Name string
}
