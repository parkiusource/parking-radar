package usecase

import (
	"time"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"golang.org/x/crypto/bcrypt"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUserName(username string) (*domain.User, error)
}

type UserUseCase struct {
	UserRepository UserRepository
}

func (u *UserUseCase) Register(username, password, email string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
		Email:        email,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	if err := u.UserRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) FindByID(id uint) (*domain.User, error) {
	return u.UserRepository.FindByID(id)
}
