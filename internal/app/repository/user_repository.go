package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

type IUserRepository interface {
	Create(user *domain.User) error
	FindByID(id uint) (*domain.User, error)
	FindByUserName(username string) (*domain.User, error)
	Update(user *domain.User) error
	Delete(id uint) error
}
