package usecase

import (
	"time"

	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/app/repository"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -source=./user_uc.go -destination=mocks/mock_user_uc.go -package=mocks
type IUserUseCase interface {
	Register(username, password, email string) (*domain.User, error)
	FindByID(id uint) (*domain.User, error)
	UpdateUser(id uint, username, email, password string) (*domain.User, error)
	DeleteUser(id uint) error
}

type UserUseCase struct {
	UserRepository repository.IUserRepository
}

func NewUserUseCase(userRepo repository.IUserRepository) IUserUseCase {
	return &UserUseCase{
		UserRepository: userRepo,
	}
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

func (u *UserUseCase) UpdateUser(id uint, username, email, password string) (*domain.User, error) {
	user, err := u.UserRepository.FindByID(id)
	if err != nil {
		return nil, err
	}

	user.Username = username
	user.Email = email
	user.UpdatedAt = time.Now()

	if password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}

	if err := u.UserRepository.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (u *UserUseCase) DeleteUser(id uint) error {
	return u.UserRepository.Delete(id)
}
