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
	Publish(topic string, ctx context.Context, payload Payload)
	Subscribe(topic string, fn Handler) error
	Unsubscribe(topic string, fn Handler) error
	Close(topic string)
}

// NewMessageBus creates in memory command bus
func NewMessageBus(maxConcurrentCalls int) MessageBus {
	return &loggableMessageBus{messagebus.New(maxConcurrentCalls)}
}

type loggableMessageBus struct {
	bus messagebus.MessageBus
}

func (b *loggableMessageBus) Publish(topic string, ctx context.Context, p Payload) {
	b.bus.Publish(topic, ctx, p)
}

func (b *loggableMessageBus) Subscribe(topic string, fn Handler) error {
	return b.bus.Subscribe(topic, fn)
}

func (b *loggableMessageBus) Unsubscribe(topic string, fn Handler) error {
	return b.bus.Unsubscribe(topic, fn)
}

func (b *loggableMessageBus) Close(topic string) {
	b.bus.Close(topic)
}
