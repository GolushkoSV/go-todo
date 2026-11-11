package main

import (
	core_logger "TODO-GO/internal/core/logger"
	core_postgres_pool "TODO-GO/internal/core/repository/postgres"
	core_http_middleware "TODO-GO/internal/core/transport/http/middleware"
	core_http_server "TODO-GO/internal/core/transport/http/server"
	users_postgres_repository "TODO-GO/internal/features/users/repository/postgres"
	users_service "TODO-GO/internal/features/users/service"
	users_transport_http "TODO-GO/internal/features/users/transport/http"
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("initializing postgres connection pool")
	// Реализуем слои Http обработчик запроса -> Сервис -> Репозиторий -> в базу
	pool, err := core_postgres_pool.NewConnectionPool(ctx, core_postgres_pool.NewConfigMust())
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing feature", zap.String("feature", "users"))

	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHttpHandler(usersService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		[]core_http_middleware.Middleware{
			core_http_middleware.RequestID(),
			core_http_middleware.Logger(logger),
			core_http_middleware.Panic(),
			core_http_middleware.Trace(),
		},
	)

	apiVersionRouter := core_http_server.NewAPIVersionRouter(core_http_server.ApiVersion1)
	apiVersionRouter.RegisterRoutes(usersTransportHTTP.Routes()...)

	httpServer.RegisterApiRoutes(apiVersionRouter)

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("http server Run %w", zap.Error(err))
	}
}
