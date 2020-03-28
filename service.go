package main

import (
	"context"

	messagebus "github.com/vardius/message-bus"
)

// Handler function
type Handler interface{}

// Payload of the message
type Payload []byte

// Service allows to subscribe/publish messages
type Service interface {
	Publish(ctx context.Context, topic string, payload Payload)
	Subscribe(topic string, fn Handler) error
	Unsubscribe(topic string, fn Handler) error
	Close(topic string)
}

type service struct {
	bus messagebus.MessageBus
}

func newService(maxConcurrentCalls int) Service {
	return &service{messagebus.New(maxConcurrentCalls)}
}

func (s *service) Publish(ctx context.Context, topic string, p Payload) {
	s.bus.Publish(topic, ctx, p)
}

func (s *service) Subscribe(topic string, fn Handler) error {
	return s.bus.Subscribe(topic, fn)
}

func (s *service) Unsubscribe(topic string, fn Handler) error {
	return s.bus.Unsubscribe(topic, fn)
}

func (s *service) Close(topic string) {
	s.bus.Close(topic)
}
