package users_service

import (
	core_domain "TODO-GO/internal/core/domain"
	"context"
	"fmt"
)

func (s *UserService) CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	// validate
	if err := user.Validate(); err != nil {
		return core_domain.User{}, fmt.Errorf("validate user domain: %w", err)
	}

	// create
	createdUser, err := s.UsersRepository.CreateUser(ctx, user)
	if err != nil {
		return core_domain.User{}, fmt.Errorf("create user: %w", err)
	}

	return createdUser, nil
}
