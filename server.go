package main

import (
	"context"
	"io"

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

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return new(empty.Empty), nil
}

// Subscribe subscribes to a topic
// Will unsubscribe when stream.Send returns error
func (s *server) Subscribe(r *pubsub_proto.SubscribeRequest, stream pubsub_proto.MessageBus_SubscribeServer) error {
	errCh := make(chan error)
	defer close(errCh)

	handler := func(ctx context.Context, payload Payload) {
		s.logger.Debug(ctx, "gRPC Server|Subscribe] %s %s", r.GetTopic(), payload)

		err := stream.Send(&pubsub_proto.SubscribeResponse{
			Payload: payload,
		})

		if err != nil {
			errCh <- err
		}
	}

	ctx := context.Background()

	s.logger.Info(ctx, "gRPC Server|Subscribe] %s", r.GetTopic())

	s.bus.Subscribe(r.GetTopic(), handler)

	err := <-errCh

	s.bus.Unsubscribe(r.GetTopic(), handler)

	if err == io.EOF {
		s.logger.Info(ctx, "gRPC Server|Unsubscribe] %s - Stream closed, no more input is available", r.GetTopic())

		return nil
	}

	s.logger.Info(ctx, "gRPC Server|Unsubscribe] %s - %s", r.GetTopic(), err.Error())

	return err
}
