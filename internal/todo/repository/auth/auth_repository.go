package auth

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	authmodel "project-evermos/internal/todo/model/auth"
)

// Repository handles data access for auth domain.
type Repository struct{
	db *gorm.DB
}

// NewRepository constructs a new Repository.
func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) FindByPhone(notelp string) (*authmodel.User, error) {
	var u authmodel.User
	if err := r.db.Where("notelp = ?", strings.TrimSpace(notelp)).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) FindByEmail(email string) (*authmodel.User, error) {
	var u authmodel.User
	if err := r.db.Where("email = ?", strings.TrimSpace(email)).First(&u).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &u, nil
}

func (r *Repository) CreateUser(u *authmodel.User) error {
	return r.db.Create(u).Error
}