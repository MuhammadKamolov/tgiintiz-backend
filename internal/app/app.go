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
	"tgiintiz/pkg/postgres"
	"time"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	pool, err := postgres.New(ctx, cfg.DB.DSN())
	if err != nil {
		return err
	}

	defer pool.Close()

	repos := repository.NewRepositories(pool)
	service := service.NewServices(repos, auth.TokenConfig{
		Secret:     cfg.JWT.Secret,
		AccessTTL:  cfg.JWT.AccessTTL,
		RefreshTTL: cfg.JWT.RefreshTTL,
	})

	handlers := handler.NewHandler(service)

	server := &http.Server{
		Addr:    cfg.Server.Port,
		Handler: handlers.InitRoutes(),
	}

	return serve(ctx, server, 10*time.Second)
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
