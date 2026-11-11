package users_postgres_repository

import (
	core_domain "TODO-GO/internal/core/domain"
	"context"
	"fmt"
)

func (r *UsersRepository) CreateUser(ctx context.Context, user core_domain.User) (core_domain.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.pool.OperationTimeout())
	defer cancel()

	query := `
	INSERT INTO users.user (full_name, phone_number)
	VALUES ($1, $2)
	RETURNING id, version, full_name, phone_number
	`

	row := r.pool.QueryRow(ctx, query, user.FullName, user.PhoneNumber)
	var userModel UserModel
	err := row.Scan(
		&userModel.ID,
		&userModel.Version,
		&userModel.FullName,
		&userModel.PhoneNumber,
	)

	if err != nil {
		return core_domain.User{}, fmt.Errorf("scan error: %w", err)
	}

	userDomin := core_domain.User{
		ID:          userModel.ID,
		Version:     userModel.Version,
		FullName:    userModel.FullName,
		PhoneNumber: userModel.PhoneNumber,
	}

	return userDomin, nil
}
