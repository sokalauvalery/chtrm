package console

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type console struct{}

// New constructs console ui struct
func New() UI {
	return console{}
}

func (console) ReadInput() (string, error) {
	reader := bufio.NewReader(os.Stdin)

	for true {
		message, err := reader.ReadString('\n')
		if err != nil {
			return "", fmt.Errorf("failed to get input %v", err)
		}

		message = strings.TrimSpace(message)
		if message != "" {
			return message, nil
		}
	}
	return "", nil
}

func (console) WriteToUser(msg string) error {
	fmt.Println(msg)
	return nil
}
