package category

import "time"

// Category maps to category table
type Category struct {
	ID           uint       `gorm:"column:id;primaryKey"`
	NamaCategory string     `gorm:"column:nama_category;size:255;not null"`
	CreatedAt    *time.Time `gorm:"column:created_at"`
	UpdatedAt    *time.Time `gorm:"column:updated_at"`
}

func (Category) TableName() string { return "category" }