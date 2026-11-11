package users_transport_http

import (
	core_domain "TODO-GO/internal/core/domain"
	core_logger "TODO-GO/internal/core/logger"
	core_http_request "TODO-GO/internal/core/transport/http/request"
	core_http_response "TODO-GO/internal/core/transport/http/response"
	"net/http"
)

type CreateUserRequestDTO struct {
	FullName    string  `json:"full_name" validate:"required,min=3,max=100"`
	PhoneNumber *string `json:"phone_number" validate:"omitempty,min=10,max=15,startswith=+"` //omitempty работают правила валидации, только если поле было передано
}

type CreateUserResponseDTO struct {
	ID          int     `json:"id"`
	Version     int     `json:"version"`
	FullName    string  `json:"full_name"`
	PhoneNumber *string `json:"phone_number"`
}

func (h *UsersHttpHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	log := core_logger.FromContext(ctx)
	responseHandler := core_http_response.NewHttpResponseHandler(log, w)

	log.Debug("invoke CreateUser handler")

	var request CreateUserRequestDTO
	if err := core_http_request.DecodeAndValidateRequest(r, &request); err != nil {
		responseHandler.ErrorResponse(err, "failed to decode and validate request")
		return
	}

	userDomain := domainFromDto(request)
	user, err := h.usersService.CreateUser(ctx, userDomain)
	if err != nil {
		responseHandler.ErrorResponse(err, "failed to create user")
		return
	}

	responseDto := dtoFromDomain(user)
	responseHandler.JsonResponse(responseDto, http.StatusCreated)
}

func domainFromDto(dto CreateUserRequestDTO) core_domain.User {
	return core_domain.NewUserUninitialized(dto.FullName, dto.PhoneNumber)
}

func dtoFromDomain(user core_domain.User) CreateUserResponseDTO {
	return CreateUserResponseDTO{
		ID:          user.ID,
		Version:     user.Version,
		FullName:    user.FullName,
		PhoneNumber: user.PhoneNumber,
	}
}
