package toko

import "time"

// Toko represents a store owned by a user.
type Toko struct {
	ID        uint      `gorm:"primaryKey;column:id"`
	IDUser    uint      `gorm:"column:id_user"`
	NamaToko  string    `gorm:"column:nama_toko"`
	UrlFoto   string    `gorm:"column:url_foto"`
	UpdatedAt time.Time `gorm:"column:updated_at"`
	CreatedAt time.Time `gorm:"column:created_at"`
}

func (Toko) TableName() string { return "toko" }