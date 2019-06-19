package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/golang/protobuf/proto"
	"github.com/golang/protobuf/ptypes/empty"
	pubsub_mock "github.com/vardius/pubsub/mock_proto"
	pubsub_proto "github.com/vardius/pubsub/proto"
)

var topic = "my-topic"
var msg = []byte("Hello you!")

var emptyResponse *empty.Empty

var subscribeResponse = &pubsub_proto.SubscribeResponse{Payload: msg}
var publishRequest = &pubsub_proto.PublishRequest{Topic: topic, Payload: msg}
var subscribeRequest = &pubsub_proto.SubscribeRequest{Topic: topic}

func TestServer(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Create mock for the stream returned by Subscribe
	stream := pubsub_mock.NewMockMessageBus_SubscribeClient(ctrl)
	// Set expectation on receiving.
	stream.EXPECT().Recv().Return(subscribeResponse, nil)
	stream.EXPECT().CloseSend().Return(nil)

	// Create mock for the client interface.
	client := pubsub_mock.NewMockMessageBusClient(ctrl)
	// Set expectation on Publish
	client.EXPECT().Publish(gomock.Any(), publishRequest).Return(emptyResponse, nil)
	// Set expectation on Subscribe
	client.EXPECT().Subscribe(gomock.Any(), subscribeRequest).Return(stream, nil)

	if err := testPubsub(client); err != nil {
		t.Fatalf("Test failed: %v", err)
	}
}

func testPubsub(client pubsub_proto.MessageBusClient) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// test Subscribe
	stream, err := client.Subscribe(ctx, subscribeRequest)
	if err != nil {
		return err
	}
	if err := stream.CloseSend(); err != nil {
		return err
	}
	got, err := stream.Recv()
	if err != nil {
		return err
	}
	if !proto.Equal(got, subscribeResponse) {
		return fmt.Errorf("stream.Recv() = %v, want %v", got, subscribeResponse)
	}

	// test Publish
	_, err = client.Publish(ctx, publishRequest)
	if err != nil {
		return err
	}

	return nil
}
