package main

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())

	if bus == nil {
		t.Fail()
	}
}

func TestSubscribePublish(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())
	ctx := context.Background()
	c := make(chan error)

	bus.Subscribe("topic", func(ctx context.Context, _ []byte) {
		c <- nil
	})

	bus.Publish("topic", ctx, []byte("ok"))

	for {
		select {
		case <-ctx.Done():
			t.Fatal(ctx.Err())
			return
		case err := <-c:
			if err != nil {
				t.Error(err)
			}
			return
		}
	}
}

func TestUnsubscribe(t *testing.T) {
	bus := NewMessageBus(runtime.NumCPU())
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	handler := func(ctx context.Context, _ []byte) {
		t.Fail()
	}

	bus.Subscribe("topic", handler)
	bus.Unsubscribe("topic", handler)

	bus.Publish("topic", ctx, []byte("ok"))

	for {
		select {
		case <-ctx.Done():
			return
		}
	}
}
