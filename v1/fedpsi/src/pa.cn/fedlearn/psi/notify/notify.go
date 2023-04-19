package notify

import (
	"context"
)

// Some events should be notified to external system

type Notification struct {
	Sender   string      `json:"sender"`
	Receiver string      `json:"receiver"`
	Payload  interface{} `json:"payload"`
}

type Notifier interface {
	Notify(ctx context.Context, msg string) error
}
