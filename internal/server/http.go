package server

import (
	"context"
	"errors"
	"net/http"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/rksouthasiait/msv_utils/logger"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
	"google.golang.org/protobuf/encoding/protojson"
)

type RegisteringHTTPHandlerServer func(mux *runtime.ServeMux)

func NewHTTPServer(
	ctx context.Context,
	waitGroup *errgroup.Group,
	address string,
	registering RegisteringHTTPHandlerServer,
) {
	jsonOption := runtime.WithMarshalerOption(
		runtime.MIMEWildcard,
		&runtime.JSONPb{
			MarshalOptions: protojson.MarshalOptions{
				UseProtoNames: true,
			},
			UnmarshalOptions: protojson.UnmarshalOptions{
				DiscardUnknown: true,
			},
		},
	)
	grpcMux := runtime.NewServeMux(jsonOption)
	registering(grpcMux)

	mux := http.NewServeMux()
	mux.Handle("/", grpcMux)

	httpServer := &http.Server{
		Handler: logger.HttpLogger(mux),
		Addr:    address,
	}

	waitGroup.Go(func() error {
		log.Info().Msgf("Start HTTP gateway server at %v", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			if errors.Is(err, http.ErrServerClosed) {
				return nil
			}
			log.Error().Err(err).Msgf("Cannot start HTTP gateway server: %v", err)
			return err
		}
		return nil
	})

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Graceful shutdown http server")

		err := httpServer.Shutdown(context.Background())
		if err != nil {
			log.Error().Err(err).Msg("Shutdown HTTP gateway server failed")
			return err
		}

		log.Info().Msg("HTTP gateway server is stopped")
		return nil
	})
}
