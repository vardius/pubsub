package main

import (
	"context"
	"io"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/vardius/golog"

	"github.com/vardius/pubsub/v2/proto"
)

type server struct {
	bus    Service
	logger golog.Logger
}

// newServer returns new pub/sub server object
func newServer(bus Service, logger golog.Logger) proto.PubSubServer {
	return &server{bus, logger}
}

// Publish publishes message payload to all topic handlers
func (s *server) Publish(ctx context.Context, r *proto.PublishRequest) (*empty.Empty, error) {
	s.logger.Debug(ctx, "Publish: %s %s", r.GetTopic(), r.GetPayload())

	s.bus.Publish(ctx, r.GetTopic(), r.GetPayload())

	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	return new(empty.Empty), nil
}

// Subscribe subscribes to a topic,
// will unsubscribe when stream.Send returns error
func (s *server) Subscribe(r *proto.SubscribeRequest, stream proto.PubSub_SubscribeServer) error {
	errCh := make(chan error)
	defer close(errCh)

	handler := func(ctx context.Context, payload Payload) {
		s.logger.Debug(ctx, "Subscribe: %s %s", r.GetTopic(), payload)

		err := stream.Send(&proto.SubscribeResponse{
			Payload: payload,
		})

		if err != nil {
			errCh <- err
		}
	}

	ctx := context.Background()

	s.logger.Info(ctx, "Subscribe: %s", r.GetTopic())

	if err := s.bus.Subscribe(r.GetTopic(), handler); err != nil {
		return err
	}

	err := <-errCh

	if err := s.bus.Unsubscribe(r.GetTopic(), handler); err != nil {
		return err
	}

	if err == io.EOF {
		s.logger.Info(ctx, "Unsubscribe: %s - Stream closed, no more input is available", r.GetTopic())

		return nil
	}

	s.logger.Info(ctx, "Unsubscribe: %s - %s", r.GetTopic(), err.Error())

	return err
}
