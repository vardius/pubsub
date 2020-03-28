package main

import (
	"context"
	"runtime"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	s := newService(runtime.NumCPU())

	if s == nil {
		t.Fail()
	}
}

func TestSubscribePublish(t *testing.T) {
	s := newService(runtime.NumCPU())
	ctx := context.Background()
	c := make(chan error)

	handler := func(ctx context.Context, _ []byte) {
		c <- nil
	}

	if err := s.Subscribe("topic", handler); err != nil {
		t.Fatal(err)
	}

	s.Publish(ctx, "topic", []byte("ok"))

	ctxDoneCh := ctx.Done()
	for {
		select {
		case <-ctxDoneCh:
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
	s := newService(runtime.NumCPU())
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	handler := func(ctx context.Context, _ []byte) {
		t.Fail()
	}

	if err := s.Subscribe("topic", handler); err != nil {
		t.Fatal(err)
	}
	if err := s.Unsubscribe("topic", handler); err != nil {
		t.Fatal(err)
	}

	s.Publish(ctx, "topic", []byte("ok"))

	ctxDoneCh := ctx.Done()
	for {
		select {
		case <-ctxDoneCh:
			return
		}
	}
}
