package core_http_middleware

import (
	core_logger "TODO-GO/internal/core/logger"
	core_http_response "TODO-GO/internal/core/transport/http/response"
	"context"
	"net/http"
	"time"
	"uuid"

	"go.uber.org/zap"
)

const requestIdHeader = "X-Request-Id"

func RequestID() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestID := r.Header.Get(requestIdHeader)
			if requestID == "" {
				requestID = uuid.New().String()
			}

			r.Header.Set(requestIdHeader, requestID)
			w.Header().Set(requestIdHeader, requestID)

			next.ServeHTTP(w, r)
		})
	}
}

func Logger(log *core_logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestId := r.Header.Get(requestIdHeader)
			logger := log.With(
				zap.String("request_id", requestId),
				zap.String("url", r.URL.String()),
			)

			ctx := context.WithValue(r.Context(), "log", logger)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Panic() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			responseHandler := core_http_response.NewHttpResponseHandler(log, w)

			defer func() {
				// recover восстановление после паники
				if p := recover(); p != nil {
					responseHandler.PanicResponse(p, "during handle HTTP request got unexpected panic")
				}
			}()

			next.ServeHTTP(w, r)
		})
	}
}

func Trace() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := r.Context()
			log := core_logger.FromContext(ctx)
			timeBefore := time.Now()

			rw := core_http_response.NewResponseWriter(w)
			log.Debug(
				">>> incoming HTTP request",
				zap.Time("time", timeBefore.UTC()),
			)

			next.ServeHTTP(rw, r)

			log.Debug(
				"<<< done HTTP response",
				zap.Int("status_code", rw.GetStatusCode()),
				zap.Duration("latency", time.Now().Sub(timeBefore)),
			)
		})
	}
}
