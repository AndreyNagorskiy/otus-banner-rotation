package main

import (
	"context"
	"log/slog"
	"os/signal"
	"syscall"

	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/algorithms"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/app"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/handlers/grpc"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/logger"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/server/grpc"
	"github.com/AndreyNagorskiy/otus-banner-rotation/internal/storage"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	cfg := MustLoad()
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

	// TODO RabbitMQ/Kafka
	rep := storage.NewRepository(dbPool)
	b := algorithms.BanditUCB1{}
	bannerRotation := app.New(l, rep, &b)

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
