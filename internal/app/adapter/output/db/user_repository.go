package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type UserRepository struct {
	DB *gorm.DB
}

func (r *UserRepository) Create(user *domain.User) error {
	return r.DB.Create(user).Error
}

func (r *UserRepository) FindByID(id uint) (*domain.User, error) {
	var user domain.User
	if err := r.DB.First(&user, id).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) FindByUserName(username string) (*domain.User, error) {
	var user domain.User
	if err := r.DB.Where("username = ?", username).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Update(user *domain.User) error {
	return r.DB.Save(user).Error
}

func (r *UserRepository) Delete(id uint) error {
	return r.DB.Delete(&domain.User{}, id).Error
}
