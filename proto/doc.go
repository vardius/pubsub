/*
Package proto contains protocol buffers for gRPC messagebus event ingestion and delivery system.
# Using gRPC client
## Publish example:
	package main

	import (
		"context"
		"fmt"
		"os"
		"time"

		"google.golang.org/grpc"
		"google.golang.org/grpc/keepalive"
		pubsub_proto "github.com/vardius/pubsub/proto"
	)

	func main() {
		host:= "0.0.0.0"
		port:= 9090
		ctx := context.Background()

		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
				Timeout:             20 * time.Second, // wait 20 second for ping ack before considering the connection dead
				PermitWithoutStream: true,             // send pings even without active streams
			}),
		}

		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), opts...)
		if err != nil {
			os.Exit(1)
		}
		defer conn.Close()

		client := pubsub_proto.NewMessageBusClient(pubsubConn)

		client.Publish(ctx, &pubsub_proto.PublishRequest{
			Topic:   "my-topic",
			Payload: []byte("Hello you!"),
		})
	}
## Subscribe example:
	package main

	import (
		"context"
		"fmt"
		"os"
		"time"

		"google.golang.org/grpc"
		"google.golang.org/grpc/keepalive"
		pubsub_proto "github.com/vardius/pubsub/proto"
	)

	func main() {
		host:= "0.0.0.0"
		port:= 9090
		ctx := context.Background()

		opts := []grpc.DialOption{
			grpc.WithInsecure(),
			grpc.WithKeepaliveParams(keepalive.ClientParameters{
				Time:                10 * time.Second, // send pings every 10 seconds if there is no activity
				Timeout:             20 * time.Second, // wait 20 second for ping ack before considering the connection dead
				PermitWithoutStream: true,             // send pings even without active streams
			}),
		}

		conn, err := grpc.DialContext(ctx, fmt.Sprintf("%s:%d", host, port), opts...)
		if err != nil {
			os.Exit(1)
		}
		defer conn.Close()

		client := pubsub_proto.NewMessageBusClient(pubsubConn)

		stream, err := client.Subscribe(ctx, &pubsub_proto.SubscribeRequest{
			Topic: "my-topic",
		})
		if err != nil {
			os.Exit(1)
		}

		for {
			resp, err := stream.Recv()
			if err != nil {
				os.Exit(1) // stream closed or error
			}

			fmt.Println(resp.GetPayload())
		}
	}
*/
package proto
