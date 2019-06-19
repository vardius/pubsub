package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	pubsub_proto "github.com/vardius/pubsub/proto"
)

type server struct {
	bus MessageBus
}

// NewServer returns new messagebus server object
func NewServer(bus MessageBus) pubsub_proto.MessageBusServer {
	return &server{bus}
}

// Publish publishes message payload to all topic handlers
func (s *server) Publish(ctx context.Context, r *pubsub_proto.PublishRequest) (*empty.Empty, error) {
	s.bus.Publish(r.GetTopic(), ctx, r.GetPayload())

	return new(empty.Empty), ctx.Err()
}

// Subscribe subscribes to a topic
func (s *server) Subscribe(r *pubsub_proto.SubscribeRequest, stream pubsub_proto.MessageBus_SubscribeServer) error {
	done := make(chan error)
	defer close(done)

	handler := func(_ context.Context, payload Payload) {
		err := stream.Send(&pubsub_proto.SubscribeResponse{
			Payload: payload,
		})

		if err != nil {
			done <- err
		}
	}

	s.bus.Subscribe(r.GetTopic(), handler)

	err := <-done

	s.bus.Unsubscribe(r.GetTopic(), handler)

	if err == nil {
		return nil
	}

	return err
}
