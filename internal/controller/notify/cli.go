package notify

import (
	"fmt"
	"io"
)

// Cli для вывода на экран
type Cli struct {
	bodyMessage string
	writer      io.Writer
}

// NewCli возвращает новый Cli для вывода на экран
func NewCli(bodyMessage string, writer io.Writer) *Cli {
	return &Cli{
		bodyMessage: bodyMessage,
		writer:      writer,
	}
}

// Send вывод на экран
func (n Cli) Send() (string, error) {
	fmt.Fprint(n.writer, n.bodyMessage)
	return "", nil
}
