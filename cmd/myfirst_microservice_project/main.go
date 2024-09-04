package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"os/signal"
	"slices"
	"syscall"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/michael/myfirst_microservice_project/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

func main() {
	configPath := "."
	if slices.Contains(os.Args, "--config") {
		configPath = "../.."
	}

	config, err := config.LoadConfig(configPath)
	if err != nil {
		//log.Fatal().Msgf("Failed to load config: %v", err)
	}

	if config.App.AppEnv == "debug" {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	ctx, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	connPool, err := pgxpool.New(ctx, config.App.DB.Source)
	if err != nil {
		var msg string
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			default:
				msg = fmt.Sprintf("Postgres error: %s", pgErr.Code)
			}
		} else {
			msg = fmt.Sprintf("Postgres database connection failed! Error: %s", err)
		}
		log.Fatal().Msgf("Failed connect to db: %v", msg)
	}
	defer connPool.Close()

	if err := connPool.Ping(ctx); err != nil {
		log.Fatal().Msgf("Failed to ping database: %v", err)
	}

	log.Info().Msgf("PostgreSQL database connected successfully! Host: %v", config.App.DB.Host)
}
