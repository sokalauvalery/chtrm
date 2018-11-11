package console

// UI interface contains methods for interaction with user
type UI interface {
	ReadInput() (string, error)
	WriteToUser(string) error
}
