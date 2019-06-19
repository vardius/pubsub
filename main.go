package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	pubsub_proto "github.com/vardius/pubsub/proto"
	"github.com/vardius/shutdown"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	grpc_health "google.golang.org/grpc/health"
	grpc_health_proto "google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/keepalive"
)

func main() {
	ctx := context.Background()
	bus := NewMessageBus(Env.QueueSize)

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(func(ctx context.Context, rec interface{}) (err error) {
			return grpc.Errorf(codes.Internal, "Recovered in f %v", rec)
		}),
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             5 * time.Minute, // If a client pings more than once every 5 minutes, terminate the connection
			PermitWithoutStream: true,            // Allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    2 * time.Hour,    // Ping the client if it is idle for 2 hours to ensure the connection is still active
			Timeout: 20 * time.Second, // Wait 20 second for the ping ack before assuming the connection is dead
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)

	pubsubServer := NewServer(bus)
	healthServer := grpc_health.NewServer()

	healthServer.SetServingStatus("pubsub", grpc_health_proto.HealthCheckResponse_SERVING)

	pubsub_proto.RegisterMessageBusServer(grpcServer, pubsubServer)
	grpc_health_proto.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Env.Host, Env.Port))
	if err != nil {
		log.Fatalf("tcp failed to listen %s:%d\n%v\n", Env.Host, Env.Port, err)
	}

	stop := func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		log.Print(ctx, "shutting down...\n")

		grpcServer.GracefulStop()
	}

	go func() {
		log.Printf("failed to serve: %v\n", grpcServer.Serve(lis))
		stop()
		os.Exit(1)
	}()

	log.Printf("tcp running at %s:%d\n", Env.Host, Env.Port)

	shutdown.GracefulStop(stop)
}
