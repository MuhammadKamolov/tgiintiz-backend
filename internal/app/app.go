package app

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"tgiintiz/internal/config"
	"tgiintiz/internal/handler"
	"tgiintiz/internal/repository"
	"tgiintiz/internal/service"
	"tgiintiz/internal/service/auth"
	"tgiintiz/pkg/logger"
	"tgiintiz/pkg/postgres"
	"tgiintiz/pkg/redis"
	"time"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	log, closeLog, err := logger.New(logger.Config{
		Env:   cfg.App.Env,
		Level: cfg.App.LogLevel,
		File:  cfg.App.LogFile,
	})
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer closeLog()
	slog.SetDefault(log)

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.New(ctx, postgres.Config{
		Host:     cfg.DB.Host,
		Port:     cfg.DB.Port,
		User:     cfg.DB.User,
		Password: cfg.DB.Password,
		Name:     cfg.DB.Name,
		SSLMode:  cfg.DB.SSLMode,
	})
	if err != nil {
		return err
	}
	defer pool.Close()

	redisClient, err := redis.New(ctx, redis.Config{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	defer func() {
		if err := redisClient.Close(); err != nil {
			slog.Error("failed to close redis", "error", err)
		}
	}()

	repos := repository.NewRepositories(pool, redisClient)
	service := service.NewServices(repos, auth.TokenConfig{
		Secret:     cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
	})

	handlers := handler.NewHandler(service)

	server := &http.Server{
		Addr:         cfg.Server.Port,
		Handler:      handlers.InitRoutes(),
		ReadTimeout:  cfg.Server.ReadTimeout,
		WriteTimeout: cfg.Server.WriteTimeout,
	}

	return serve(ctx, server, cfg.Server.ShutdownTimeout)
}

func serve(ctx context.Context, server *http.Server, shutdownTimeout time.Duration) error {
	serverErrors := make(chan error, 1)
	go func() {
		slog.Info("server started", "addr", server.Addr)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			serverErrors <- fmt.Errorf("server error: %w", err)
		}
	}()

	select {
	case err := <-serverErrors:
		return fmt.Errorf("fatal server error: %w", err)
	case <-ctx.Done():
	}

	slog.Info("shutting down")
	shutdownCtx, cancel := context.WithTimeout(context.Background(), shutdownTimeout)
	defer cancel()
	if err := server.Shutdown(shutdownCtx); err != nil {
		server.Close()
		return fmt.Errorf("could not stop server gracefully: %w", err)
	}

	slog.Info("server stopped gracefully")

	return nil
}
