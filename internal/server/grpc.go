package server

import (
	"context"
	"errors"
	"net"

	"github.com/rksouthasiait/msv_utils/logger"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

type RegisteringService func(grpcServer *grpc.Server, unaryInterceptor grpc.ServerOption)

func NewGRPCServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	address string,
	registering RegisteringService,
) {
	if registering == nil {
		log.Fatal().Msg("Unable to register GRPC server")
	}

	unaryInterceptor := grpc.UnaryInterceptor(logger.GprcLogger)
	grpcServer := grpc.NewServer(unaryInterceptor)
	registering(grpcServer, unaryInterceptor)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatal().Msgf("Unable create listener: %v", err)
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("Start gRPC server at %v", listener.Addr().String())
		err = grpcServer.Serve(listener)
		if err != nil {
			if errors.Is(err, grpc.ErrServerStopped) {
				return nil
			}
			log.Error().Err(err).Msgf("Cannot start gRPC server: %v", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Graceful shutdown gRPC server")

		grpcServer.GracefulStop()
		log.Info().Msg("gRPC server is stopped")
		return nil
	})
}
