package users

import "time"

// User maps to users table using existing column names
type User struct {
    ID            uint       `gorm:"column:id;primaryKey"`
    Nama          string     `gorm:"column:nama;size:255;not null"`
    KataSandi     string     `gorm:"column:kata_sandi;size:255;not null"`
    NoTelp        string     `gorm:"column:notelp;size:255;uniqueIndex;not null"`
    TanggalLahir  *time.Time `gorm:"column:tanggal lahir"`
    JenisKelamin  *string    `gorm:"column:jenis kelamin"`
    Tentang       *string    `gorm:"column:tentang"`
    Pekerjaan     string     `gorm:"column:pekerjaan;size:255;not null"`
    Email         string     `gorm:"column:email;size:255;uniqueIndex;not null"`
    IDProvinsi    string     `gorm:"column:id_provinsi;size:255;not null"`
    IDKota        string     `gorm:"column:id_kota;size:255;not null"`
    IsAdmin       *bool      `gorm:"column:isAdmin"`
    UpdatedAt     *time.Time `gorm:"column:updated_at"`
    CreatedAt     *time.Time `gorm:"column:created_at"`
}

func (User) TableName() string { return "users" }

// Alamat maps to alamat table using existing column names (note spaces in some columns)
type Alamat struct {
    ID            uint       `gorm:"column:id;primaryKey"`
    IDUser        uint       `gorm:"column:id_user;not null"`
    JudulAlamat   string     `gorm:"column:judul alamat;size:255;not null"`
    NamaPenerima  string     `gorm:"column:nama penerima;size:255;not null"`
    NoTelp        string     `gorm:"column:no telp;size:255;not null"`
    DetailAlamat  string     `gorm:"column:detail_alamat;size:255;not null"`
    UpdatedAt     *time.Time `gorm:"column:updated_at"`
    CreatedAt     *time.Time `gorm:"column:created_at"`
}

func (Alamat) TableName() string { return "alamat" }