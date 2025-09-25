package product

import (
    "time"
)

type Product struct {
    ID            uint      `gorm:"primaryKey;column:id"`
    NamaProduk    string    `gorm:"column:nama_produk"`
    Slug          string    `gorm:"column:slug"`
    HargaReseller string    `gorm:"column:harga reseller"`
    HargaKonsumen string    `gorm:"column:harga konsumen"`
    Stok          int       `gorm:"column:stok"`
    Deskripsi     string    `gorm:"column:deskripsi"`
    CreatedAt     time.Time `gorm:"column:created_at"`
    UpdatedAt     time.Time `gorm:"column:updated_at"`
    IDToko        uint      `gorm:"column:id_toko"`
    IDCategory    uint      `gorm:"column:id_category"`

    Toko     *TokoRef     `gorm:"foreignKey:IDToko;references:ID"`
    Category *CategoryRef `gorm:"foreignKey:IDCategory;references:ID"`
    Photos   []Photo      `gorm:"foreignKey:IDProduk;references:ID"`
}

func (Product) TableName() string { return "produk" }

type Photo struct {
    ID        uint      `gorm:"primaryKey;column:id"`
    IDProduk  uint      `gorm:"column:id_produk"`
    URL       string    `gorm:"column:url"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
    CreatedAt time.Time `gorm:"column:created_at"`
}

func (Photo) TableName() string { return "foto_produk" }

type CategoryRef struct {
    ID           uint      `gorm:"primaryKey;column:id"`
    NamaCategory string    `gorm:"column:nama_category"`
    CreatedAt    time.Time `gorm:"column:created_at"`
    UpdatedAt    time.Time `gorm:"column:updated_at"`
}

func (CategoryRef) TableName() string { return "category" }

type TokoRef struct {
    ID        uint      `gorm:"primaryKey;column:id"`
    NamaToko  string    `gorm:"column:nama_toko"`
    UrlFoto   string    `gorm:"column:url_foto"`
    UpdatedAt time.Time `gorm:"column:updated_at"`
    CreatedAt time.Time `gorm:"column:created_at"`
}

func (TokoRef) TableName() string { return "toko" }