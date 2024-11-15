package db

import (
	"github.com/CamiloLeonP/parking-radar/internal/app/domain"
	"gorm.io/gorm"
)

type AdminRepositoryImpl struct {
	DB *gorm.DB
}

// Create adds a new admin to the database.
func (r *AdminRepositoryImpl) Create(admin *domain.Admin) error {
	return r.DB.Create(admin).Error
}

// ExistsByAuth0UUID checks if an admin with the given Auth0 UUID exists.
func (r *AdminRepositoryImpl) ExistsByAuth0UUID(auth0UUID string) (bool, error) {
	var count int64
	if err := r.DB.Model(&domain.Admin{}).Where("auth0_uuid = ?", auth0UUID).Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByAuth0UUID retrieves an admin by its Auth0 UUID.
func (r *AdminRepositoryImpl) FindByAuth0UUID(auth0UUID string) (*domain.Admin, error) {
	var admin domain.Admin
	if err := r.DB.Where("auth0_uuid = ?", auth0UUID).First(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}

// Update modifies an existing admin.
func (r *AdminRepositoryImpl) Update(admin *domain.Admin) error {
	return r.DB.Save(admin).Error
}
