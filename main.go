package main

import (
	"context"
	"fmt"
	"net"
	"os"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	grpc_recovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	"github.com/vardius/golog"
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

	logger := golog.NewConsoleLogger(Env.Verbose)
	bus := NewMessageBus(Env.QueueSize)

	opts := []grpc_recovery.Option{
		grpc_recovery.WithRecoveryHandlerContext(func(ctx context.Context, rec interface{}) (err error) {
			logger.Error(ctx, "Recovered in f %v", rec)
			return grpc.Errorf(codes.Internal, "Recovered from panic")
		}),
	}

	grpcServer := grpc.NewServer(
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             Env.KeepaliveEnforcementPolicyMinTime,
			PermitWithoutStream: true, // Allow pings even when there are no active streams
		}),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    Env.KeepaliveParamsTime,
			Timeout: Env.KeepaliveParamsTimeout,
		}),
		grpc_middleware.WithUnaryServerChain(
			grpc_recovery.UnaryServerInterceptor(opts...),
		),
		grpc_middleware.WithStreamServerChain(
			grpc_recovery.StreamServerInterceptor(opts...),
		),
	)

	pubsubServer := NewServer(bus, logger)
	healthServer := grpc_health.NewServer()

	healthServer.SetServingStatus("pubsub", grpc_health_proto.HealthCheckResponse_SERVING)

	pubsub_proto.RegisterMessageBusServer(grpcServer, pubsubServer)
	grpc_health_proto.RegisterHealthServer(grpcServer, healthServer)

	lis, err := net.Listen("tcp", fmt.Sprintf("%s:%d", Env.Host, Env.Port))
	if err != nil {
		logger.Critical(ctx, "tcp failed to listen %s:%d\n%v\n", Env.Host, Env.Port, err)
		os.Exit(1)
	}

	stop := func() {
		ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		logger.Info(ctx, "shutting down...\n")

		grpcServer.GracefulStop()
	}

	go func() {
		logger.Critical(ctx, "failed to serve: %v\n", grpcServer.Serve(lis))
		stop()
		os.Exit(1)
	}()

	logger.Info(ctx, "tcp running at %s:%d\n", Env.Host, Env.Port)

	shutdown.GracefulStop(stop)
}
