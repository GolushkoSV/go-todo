package core_domain

import (
	core_errors "TODO-GO/internal/core/errors"
	"fmt"
	"regexp"
)

type User struct {
	ID      int
	Version int

	FullName    string  //Обязателен
	PhoneNumber *string //опционален
}

func NewUser(id int, version int, fullName string, phoneNumber *string) User {
	return User{
		ID:          id,
		Version:     version,
		FullName:    fullName,
		PhoneNumber: phoneNumber,
	}
}

func NewUserUninitialized(
	fullName string,
	phoneNumber *string,
) User {
	return NewUser(
		UninitializedID,
		UninitializedVersion,
		fullName,
		phoneNumber,
	)
}

func (u *User) Validate() error {
	fullNameLength := len([]rune(u.FullName)) // массив рун = массив символов в строке
	if fullNameLength < 3 || fullNameLength > 100 {
		return fmt.Errorf("invalid `FullName` len: %d: %w", fullNameLength, core_errors.ErrInvalidArgument)
	}

	if u.PhoneNumber != nil {
		phoneNumberLength := len([]rune(*u.PhoneNumber))
		if phoneNumberLength < 10 || phoneNumberLength > 15 {
			return fmt.Errorf("invalid `PhoneNumber` len: %d: %w", phoneNumberLength, core_errors.ErrInvalidArgument)
		}

		reg := regexp.MustCompile(`^\+[0-9]+$`)
		if !reg.MatchString(*u.PhoneNumber) {
			return fmt.Errorf("invalid `PhoneNumber` format: %w", core_errors.ErrInvalidArgument)
		}
	}

	return nil
}
