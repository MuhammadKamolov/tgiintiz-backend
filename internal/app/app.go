package app

import (
	"context"
	"fmt"
	"net/http"
	"tgiintiz/internal/config"
	v1 "tgiintiz/internal/handler/v1"
	"tgiintiz/internal/repository"
	"tgiintiz/internal/service"
	"tgiintiz/internal/service/auth"
	"tgiintiz/pkg/postgres"
)

func Run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx := context.Background()

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

	apiV1 := v1.NewHandler(service)
	routes := apiV1.InitRoutes()

	fmt.Printf("HTTP server listening on %s\n", cfg.Server.Port)
	if err := http.ListenAndServe(cfg.Server.Port, routes); err != nil {
		return fmt.Errorf("server failed to start: %w", err)
	}

	return nil
}
