package main

import (
	"context"

	messagebus "github.com/vardius/message-bus"
)

// Handler function
type Handler interface{}

// Payload of the message
type Payload []byte

// MessageBus allows to subscribe/dispatch messages
type MessageBus interface {
	Publish(ctx context.Context, topic string, payload Payload)
	Subscribe(topic string, fn Handler) error
	Unsubscribe(topic string, fn Handler) error
	Close(topic string)
}

// NewMessageBus creates in memory command bus
func NewMessageBus(maxConcurrentCalls int) MessageBus {
	return &bus{messagebus.New(maxConcurrentCalls)}
}

type bus struct {
	bus messagebus.MessageBus
}

func (b *bus) Publish(ctx context.Context, topic string, p Payload) {
	b.bus.Publish(topic, ctx, p)
}

func (b *bus) Subscribe(topic string, fn Handler) error {
	return b.bus.Subscribe(topic, fn)
}

func (b *bus) Unsubscribe(topic string, fn Handler) error {
	return b.bus.Unsubscribe(topic, fn)
}

func (b *bus) Close(topic string) {
	b.bus.Close(topic)
}
