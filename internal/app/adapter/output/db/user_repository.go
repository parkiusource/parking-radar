package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"github.com/CamiloLeonP/parking-radar/internal/db"
)

type UserRepository struct{}

func (r *UserRepository) Create(user *domain.User) error {
	return db.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := db.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) FindByUserName(username string) (*domain.User, error) {
	var user domain.User
	if err := db.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
