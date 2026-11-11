package core_http_response

import (
	core_errors "TODO-GO/internal/core/errors"
	core_logger "TODO-GO/internal/core/logger"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

type HttpResponseHandler struct {
	log *core_logger.Logger
	rw  http.ResponseWriter
}

func NewHttpResponseHandler(log *core_logger.Logger, rw http.ResponseWriter) *HttpResponseHandler {
	return &HttpResponseHandler{
		log: log,
		rw:  rw,
	}
}

func (h *HttpResponseHandler) JsonResponse(body any, statusCode int) {
	h.rw.WriteHeader(statusCode)

	if err := json.NewEncoder(h.rw).Encode(body); err != nil {
		h.log.Error("write HTTP response", zap.Error(err))
	}
}

func (h *HttpResponseHandler) ErrorResponse(err error, msg string) {
	var (
		statusCode int
		logFunc    func(string, ...zap.Field)
	)

	switch {
	case errors.Is(err, core_errors.ErrInvalidArgument):
		statusCode = http.StatusBadRequest
		logFunc = h.log.Warn
	case errors.Is(err, core_errors.ErrNotFound):
		statusCode = http.StatusNotFound
		logFunc = h.log.Debug
	case errors.Is(err, core_errors.ErrConflict):
		statusCode = http.StatusConflict
		logFunc = h.log.Warn
	default:
		statusCode = http.StatusInternalServerError
		logFunc = h.log.Error
	}

	logFunc(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

// todo Надо бы закрыть печать ошибок на фронт для Паники
func (h *HttpResponseHandler) PanicResponse(p any, msg string) {
	statusCode := http.StatusInternalServerError
	err := fmt.Errorf("unexpected panic: %v", p) // %v все что придет, преобразует в строку

	h.log.Error(msg, zap.Error(err))

	h.errorResponse(
		statusCode,
		err,
		msg,
	)
}

func (h *HttpResponseHandler) errorResponse(statusCode int, err error, msg string) {
	response := map[string]string{
		"message": msg,
		"error":   err.Error(),
	}

	h.JsonResponse(response, statusCode)
}
