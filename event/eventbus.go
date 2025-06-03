package event

import "context"

// EventBus is an interface for publishing and subscribing to events
type EventBus interface {
	Publish(ctx context.Context, topic string, payload any) error
	Subscribe(ctx context.Context, topic string, handler EventHandler) error
}

// EventHandler is an interface for handling events
type EventHandler interface {
	Handle(ctx context.Context, message any) error
}
