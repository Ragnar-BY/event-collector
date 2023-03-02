package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Ragnar-BY/event-collector/internal/config"
	"github.com/Ragnar-BY/event-collector/internal/controllers/rest"
	"github.com/Ragnar-BY/event-collector/internal/repository/clickhouse"
	"github.com/Ragnar-BY/event-collector/internal/service"
	"github.com/Ragnar-BY/event-collector/internal/usecase"

	"go.uber.org/zap"
)

func main() {
	logger, _ := zap.NewProduction()
	cfg, err := config.LoadConfig(".env")
	if err != nil {
		logger.Fatal("can not load config", zap.Error(err))
	}

	client, err := clickhouse.NewClickhouseClient(clickhouse.ClickhouseSettings{
		Addr:     cfg.ClickhouseAddr,
		Database: cfg.ClickhouseDatabase,
		Username: cfg.ClickhouseUser,
		Password: cfg.ClickhousePassword,
	})
	if err != nil {
		logger.Fatal("can not run clickhouse", zap.Error(err))
	}

	eventService := service.NewRepositoryService(client)
	eventUsecase := usecase.NewEventUsecase(eventService, logger, usecase.Config{
		NumberOfThreads: cfg.NumberOfInsertThreads,
		ChannelCapacity: cfg.ChannelCapacity,
	})
	srv := rest.NewServer(cfg.ServerAddress, logger, eventUsecase)

	go func() {
		err = srv.Run()
		if err != nil {
			logger.Error("can not run server", zap.Error(err))
		}
	}()

	// Graceful shutdown
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatal("Server forced to shutdown: ", zap.Error(err))
	}
	if err := eventService.CloseDB(); err != nil {
		logger.Fatal("Database service forced to shutdown: ", zap.Error(err))
	}

	logger.Info("Server exiting")

}
