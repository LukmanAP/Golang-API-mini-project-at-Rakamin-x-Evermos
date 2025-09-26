package category

import (
	"errors"
	"strings"

	"gorm.io/gorm"
	model "project-evermos/internal/todo/model/category"
)

var (
	ErrDuplicate = errors.New("duplicate")
	ErrNotFound  = errors.New("not found")
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) List() ([]model.Category, error) {
	var rows []model.Category
	if err := r.db.Find(&rows).Error; err != nil { return nil, err }
	return rows, nil
}

func (r *Repository) GetByID(id uint) (*model.Category, error) {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return nil, ErrNotFound }
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ExistsByName(name string) (bool, error) {
	t := strings.TrimSpace(name)
	if t == "" { return false, nil }
	var cnt int64
	if err := r.db.Model(&model.Category{}).Where("LOWER(nama_category) = ?", strings.ToLower(t)).Count(&cnt).Error; err != nil { return false, err }
	return cnt > 0, nil
}

func (r *Repository) Create(name string) (uint, error) {
	// unique check
	exists, err := r.ExistsByName(name)
	if err != nil { return 0, err }
	if exists { return 0, ErrDuplicate }
	row := model.Category{NamaCategory: strings.TrimSpace(name)}
	if err := r.db.Create(&row).Error; err != nil { return 0, err }
	return row.ID, nil
}

func (r *Repository) Update(id uint, name string) error {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return ErrNotFound }
		return err
	}
	// unique check
	exists, err := r.ExistsByName(name)
	if err != nil { return err }
	if exists && strings.ToLower(c.NamaCategory) != strings.ToLower(strings.TrimSpace(name)) {
		return ErrDuplicate
	}
	c.NamaCategory = strings.TrimSpace(name)
	return r.db.Save(&c).Error
}

func (r *Repository) Delete(id uint) error {
	var c model.Category
	if err := r.db.First(&c, id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) { return ErrNotFound }
		return err
	}
	return r.db.Delete(&c).Error
}