package repository

import "github.com/CamiloLeonP/parking-radar/internal/app/domain"

//go:generate mockgen -source=./admin_repository.go -destination=./../../test/shared/mocks/mock_admin_repository.go -package=mockgen
type IAdminRepository interface {
	Create(admin *domain.Admin) error
	ExistsByAuth0UUID(auth0UUID string) (bool, error)
	FindByAuth0UUID(auth0UUID string) (*domain.Admin, error)
	Update(admin *domain.Admin) error
}
