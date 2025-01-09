package grpc

import (
	"context"
	"net"
	"time"

	grpc_middleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/pkg/errors"
	"github.com/tguankheng016/go-ecommerce-microservice/internal/pkg/logging"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
)

const (
	maxConnectionIdle = 5
	gRPCTimeout       = 15
	maxConnectionAge  = 5
	gRPCTime          = 10
)

type GrpcServer struct {
	Grpc   *grpc.Server
	Config *GrpcOptions
}

func NewGrpcServer(config *GrpcOptions) *GrpcServer {
	if !config.Enabled {
		return nil
	}

	unaryServerInterceptors := []grpc.UnaryServerInterceptor{
		otelgrpc.UnaryServerInterceptor(),
	}
	streamServerInterceptors := []grpc.StreamServerInterceptor{
		otelgrpc.StreamServerInterceptor(),
	}

	s := grpc.NewServer(
		grpc.KeepaliveParams(keepalive.ServerParameters{
			MaxConnectionIdle: maxConnectionIdle * time.Minute,
			Timeout:           gRPCTimeout * time.Second,
			MaxConnectionAge:  maxConnectionAge * time.Minute,
			Time:              gRPCTime * time.Minute,
		}),
		grpc.StreamInterceptor(
			grpc_middleware.ChainStreamServer(
				streamServerInterceptors...,
			),
		),
		grpc.UnaryInterceptor(
			grpc_middleware.ChainUnaryServer(
				unaryServerInterceptors...,
			),
		),
	)

	return &GrpcServer{Grpc: s, Config: config}
}

func (s *GrpcServer) RunGrpcServer(configGrpc ...func(grpcServer *grpc.Server)) error {
	listen, err := net.Listen("tcp", s.Config.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	if s.Config.Development {
		reflection.Register(s.Grpc)
	}

	if len(configGrpc) > 0 {
		grpcFunc := configGrpc[0]
		if grpcFunc != nil {
			grpcFunc(s.Grpc)
		}
	}

	logging.Logger.Info("grpc server is listening on port: " + s.Config.Port)

	err = s.Grpc.Serve(listen)

	if err != nil {
		logging.Logger.Error("[grpcServer_RunGrpcServer.Serve] grpc server serve error:", zap.Error(err))
	}

	return err
}

func (s *GrpcServer) shutdown() {
	s.Grpc.Stop()
	s.Grpc.GracefulStop()
}

func RunServers(lc fx.Lifecycle, grpcServer *GrpcServer, clientFactory *GrpcClientFactory) error {
	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			if grpcServer == nil || !grpcServer.Config.Enabled {
				return nil
			}

			go func() {
				if err := grpcServer.RunGrpcServer(); err != nil {
					logging.Logger.Error("error running grpc server", zap.Error(err))
				}
			}()

			return nil
		},
		OnStop: func(ctx context.Context) error {
			if grpcServer != nil && grpcServer.Config.Enabled {
				logging.Logger.Info("closing grpc servers...")
				grpcServer.shutdown()
				logging.Logger.Info("all grpc servers shutdown gracefully...")
			}
			clientFactory.RemoveClients()
			return nil
		},
	})

	return nil
}
