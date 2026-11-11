package users_transport_http

import (
	"TODO-GO/internal/core/domain"
	core_http_server "TODO-GO/internal/core/transport/http/server"
	"context"
	"net/http"
)

type UsersHttpHandler struct {
	usersService UsersService
}

type UsersService interface {
	CreateUser(
		ctx context.Context,
		user core_domain.User,
	) (core_domain.User, error)
}

func NewUsersHttpHandler(usersService UsersService) *UsersHttpHandler {
	return &UsersHttpHandler{
		usersService: usersService,
	}
}

func (h *UsersHttpHandler) Routes() []core_http_server.Route {
	return []core_http_server.Route{
		{
			Method:  http.MethodPost,
			Path:    "/users",
			Handler: h.CreateUser,
		},
	}
}
