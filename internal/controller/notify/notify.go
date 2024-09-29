package notify

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
