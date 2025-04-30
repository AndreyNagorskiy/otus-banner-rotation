package main

import (
	"context"
	"flag"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/app"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/handlers/grpc"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/logger"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/server/grpc"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

var configFile string

func main() {
	flag.StringVar(&configFile, "config", "configs/config.yaml", "Path to configuration file")
	flag.Parse()

	cfg := MustLoad(configFile)
	l := logger.NewLogger(cfg.LogLevel)

	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	dbConnectionString := cfg.MakeDBConnectionString()
	err := storage.Migrate(dbConnectionString, false)
	if err != nil {
		l.Error("Unable to migrate database", slog.String("error", err.Error()))
		return
	}

	dbPool, err := pgxpool.New(ctx, dbConnectionString)
	if err != nil {
		l.Error("Unable to connect to database", slog.String("error", err.Error()))
		return
	}
	defer dbPool.Close()

	// TODO bandit algorithm and RabbitMQ/Kafka
	rep := storage.NewRepository(dbPool)
	bannerRotation := app.New(l, rep)

	grpcHandler := grpchandler.NewBannerRotationHandler(bannerRotation)
	grpcServer := internalgrpc.NewServer(l, grpcHandler)

	done := make(chan struct{})

	go func() {
		if err := grpcServer.Run(cfg.MakeGRPCAddr()); err != nil {
			l.Error("Grpc server failed", slog.String("error", err.Error()))
		}
		close(done)
	}()

	<-ctx.Done()
	l.Info("Shutdown signal received")

	grpcServer.Stop()
	<-done

	l.Info("Application stopped gracefully")
}
