package core_http_request

import (
	core_errors "TODO-GO/internal/core/errors"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-playground/validator/v10"
)

var requestValidator = validator.New()

func DecodeAndValidateRequest(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		return fmt.Errorf(
			"decode json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	if err := requestValidator.Struct(dest); err != nil {
		return fmt.Errorf(
			"validate json: %v: %w",
			err,
			core_errors.ErrInvalidArgument,
		)
	}

	return nil
}
