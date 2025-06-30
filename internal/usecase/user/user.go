package user

import (
	"errors"
	"fmt"
	"time"

	"github.com/Noviiich/golang-new-relic-sentry-prometheus/internal/domain"
)

type UserRepository interface {
	CreateUser(user domain.User) (domain.User, error)
	GetUserByID(id int) (domain.User, error)
	UpdateUser(user domain.User) (domain.User, error)
	DeleteUser(id int) error
}

type UserUseCase struct {
	repo UserRepository
}

func NewUserUseCase(repo UserRepository) *UserUseCase {
	return &UserUseCase{repo: repo}
}

func (uc *UserUseCase) CreateUser(user domain.User) (domain.User, error) {
	const op = "user.usecase.CreateUser"

	user.CreatedDate = time.Now()
	if user.Name == "" {
		return domain.User{}, fmt.Errorf("%s: %w", op, errors.New("name is required"))
	}

	createdUser, err := uc.repo.CreateUser(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return createdUser, nil
}

func (uc *UserUseCase) GetUserByID(id int) (domain.User, error) {
	const op = "user.usecase.GetUserByID"

	user, err := uc.repo.GetUserByID(id)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return user, nil
}

func (uc *UserUseCase) UpdateUser(user domain.User) (domain.User, error) {
	const op = "user.usecase.UpdateUser"

	updatedUser, err := uc.repo.UpdateUser(user)
	if err != nil {
		return domain.User{}, fmt.Errorf("%s: %w", op, err)
	}
	return updatedUser, nil
}

func (uc *UserUseCase) DeleteUser(id int) error {
	const op = "user.usecase.DeleteUser"

	err := uc.repo.DeleteUser(id)
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}
	return nil
}
