package toko

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	model "project-evermos/internal/todo/model/toko"
)

// Repository handles data access for toko domain.
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Create(t *model.Toko) error {
	return r.db.Create(t).Error
}

func (r *Repository) FindByUserID(userID uint) (*model.Toko, error) {
	var t model.Toko
	if err := r.db.Where("id_user = ?", userID).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) FindByID(id uint) (*model.Toko, error) {
	var t model.Toko
	if err := r.db.Where("id = ?", id).First(&t).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &t, nil
}

func (r *Repository) Update(t *model.Toko) error {
	return r.db.Save(t).Error
}

func (r *Repository) List(limit, page int, name string) ([]model.Toko, int64, error) {
	var items []model.Toko
	var count int64
	q := r.db.Model(&model.Toko{})
	if strings.TrimSpace(name) != "" {
		q = q.Where("nama_toko LIKE ?", "%"+strings.TrimSpace(name)+"%")
	}
	if err := q.Count(&count).Error; err != nil {
		return nil, 0, err
	}
	offset := (page - 1) * limit
	if err := q.Order("id DESC").Limit(limit).Offset(offset).Find(&items).Error; err != nil {
		return nil, 0, err
	}
	return items, count, nil
}