package core_http_server

import (
	core_logger "TODO-GO/internal/core/logger"
	core_http_middleware "TODO-GO/internal/core/transport/http/middleware"
	"context"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HTTPServer struct {
	mux    *http.ServeMux
	config Config
	log    *core_logger.Logger

	middlewares []core_http_middleware.Middleware
}

func NewHTTPServer(config Config, log *core_logger.Logger, middleware []core_http_middleware.Middleware) *HTTPServer {
	return &HTTPServer{
		mux:         http.NewServeMux(),
		config:      config,
		log:         log,
		middlewares: middleware,
	}
}

func (srv *HTTPServer) RegisterApiRoutes(routes ...*APIVersionRouter) {
	for _, router := range routes {
		prefix := "/api/" + string(router.apiVersion)
		srv.mux.Handle(
			prefix+"/",
			http.StripPrefix(prefix, router),
		)
	}
}

func (srv *HTTPServer) Run(ctx context.Context) error {
	mux := core_http_middleware.ChainMiddleware(srv.mux, srv.middlewares...)
	server := &http.Server{
		Addr:    srv.config.Addr,
		Handler: mux,
	}

	ch := make(chan error, 1)
	go func() {
		defer close(ch)
		srv.log.Warn("start HTTP server", zap.String("addr", srv.config.Addr))
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			ch <- err
		}
	}()

	select {
	case err := <-ch:
		if err != nil {
			return fmt.Errorf("listen and server HTTP: %w", err)
		}
	case <-ctx.Done():
		srv.log.Warn("shutdown HTTP server")
		shutdownCtx, cancel := context.WithTimeout(
			context.Background(),
			srv.config.ShutdownTimeout,
		)

		defer cancel()
		if err := server.Shutdown(shutdownCtx); err != nil {
			_ = server.Close()
			return fmt.Errorf("shutdown HTTP server: %w", err)
		}

		srv.log.Warn("HTTP server stopped")
	}

	return nil
}
