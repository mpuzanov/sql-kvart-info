package notify

import (
	"fmt"
	"io"
	"kvart-info/pkg/email"
)

// Notifier ...
type Notifier interface {
	Send() (string, error)
}

// Notify ...
type Notify struct {
	Notifier
}

// New  создание объекта
func New(obj Notifier) Notify {
	return Notify{obj}
}

// NotifyEmail для отправки по почте
type NotifyEmail struct {
	Cfg         email.Config
	BodyMessage string
	Title       string
	ToSendEmail string
}

// Send отправка по почте
func (n NotifyEmail) Send() (string, error) {
	objEmail := n.Cfg
	return objEmail.Send(n.BodyMessage, n.Title, n.ToSendEmail, "")
}

// NotifyCli для вывода на экра
type NotifyCli struct {
	BodyMessage string
	Writer      io.Writer
}

// Send вывод на экран
func (n NotifyCli) Send() (string, error) {
	fmt.Fprint(n.Writer, n.BodyMessage)
	return "", nil
}
