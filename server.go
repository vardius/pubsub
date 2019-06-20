package main

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vardius/golog"
	pubsub_proto "github.com/vardius/pubsub/proto"
)

type server struct {
	bus    MessageBus
	logger golog.Logger
}

// NewServer returns new messagebus server object
func NewServer(bus MessageBus, logger golog.Logger) pubsub_proto.MessageBusServer {
	return &server{bus, logger}
}

// Publish publishes message payload to all topic handlers
func (s *server) Publish(ctx context.Context, r *pubsub_proto.PublishRequest) (*empty.Empty, error) {
	s.logger.Debug(ctx, "gRPC Server|Publish] %s %s", r.GetTopic(), r.GetPayload())

	s.bus.Publish(ctx, r.GetTopic(), r.GetPayload())

	return new(empty.Empty), ctx.Err()
}

// Subscribe subscribes to a topic
// Will unsubscribe when stream.Send returns error
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

	ctx := context.Background()

	s.logger.Info(ctx, "gRPC Server|Subscribe] %s", r.GetTopic())

	s.bus.Subscribe(r.GetTopic(), handler)

	err := <-done

	s.logger.Info(ctx, "gRPC Server|Unsubscribe] %s %v", r.GetTopic(), err)

	s.bus.Unsubscribe(r.GetTopic(), handler)

	return err
}
