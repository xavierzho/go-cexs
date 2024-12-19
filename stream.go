package cexconns

import (
	"context"
)

type StreamClient interface {
	// Connect dial to server
	Connect(ctx context.Context, url string) error
	// Close connection
	Close() error
	// SendMessage send json payload
	SendMessage(payload any) error
	// Reconnect reset connection
	Reconnect(ctx context.Context, url string) error
	// Ping heartbreak
	Ping(ctx context.Context) error
}

type PublicSubscriber interface {
	StreamClient
	Subscribe(ctx context.Context)
	UnSubscribe(ctx context.Context)
}

type PrivateSubscriber interface {
	StreamClient
	Login() error
}
