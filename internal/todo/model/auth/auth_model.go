package auth

import "time"

// User represents the users table.
// Using gorm tags to map to existing columns with spaces and different names.
type User struct {
    ID           uint      `gorm:"column:id;primaryKey"`
    Nama         string    `gorm:"column:nama;size:255;not null"`
    KataSandi    string    `gorm:"column:kata_sandi;size:255;not null"`
    NoTelp       string    `gorm:"column:notelp;size:255;uniqueIndex;not null"`
    TanggalLahir *time.Time `gorm:"column:tanggal lahir"`
    JenisKelamin *string   `gorm:"column:jenis kelamin"`
    Tentang      *string   `gorm:"column:tentang"`
    Pekerjaan    string    `gorm:"column:pekerjaan;size:255;not null"`
    Email        string    `gorm:"column:email;size:255;uniqueIndex;not null"`
    IDProvinsi   string    `gorm:"column:id_provinsi;size:255;not null"`
    IDKota       string    `gorm:"column:id_kota;size:255;not null"`
    IsAdmin      *bool     `gorm:"column:isAdmin"`
    UpdatedAt    *time.Time `gorm:"column:updated_at"`
    CreatedAt    *time.Time `gorm:"column:created_at"`
}

func (User) TableName() string { return "users" }