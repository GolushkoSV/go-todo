package users_service

import (
	core_domain "TODO-GO/internal/core/domain"
	"context"
)

type UserService struct {
	UsersRepository UsersRepository
}

type UsersRepository interface {
	CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error)
}

func NewUsersService(repository UsersRepository) *UserService {
	return &UserService{
		UsersRepository: repository,
	}
}
