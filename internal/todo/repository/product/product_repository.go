package product

import (
    "errors"
    "strings"

    prodmodel "project-evermos/internal/todo/model/product"

    "gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// List products with filters and pagination
func (r *Repository) List(filter ListFilter) ([]prodmodel.Product, int64, error) {
    var items []prodmodel.Product
    var count int64

    q := r.db.Model(&prodmodel.Product{})

    if s := strings.TrimSpace(filter.NamaProduk); s != "" {
        q = q.Where("nama_produk LIKE ?", "%"+s+"%")
    }
    if filter.CategoryID > 0 {
        q = q.Where("id_category = ?", filter.CategoryID)
    }
    if filter.TokoID > 0 {
        q = q.Where("id_toko = ?", filter.TokoID)
    }
    if filter.MinHarga != nil {
        q = q.Where("CAST(`harga konsumen` AS SIGNED) >= ?", *filter.MinHarga)
    }
    if filter.MaxHarga != nil {
        q = q.Where("CAST(`harga konsumen` AS SIGNED) <= ?", *filter.MaxHarga)
    }

    if err := q.Count(&count).Error; err != nil {
        return nil, 0, err
    }

    offset := (filter.Page - 1) * filter.Limit
    if offset < 0 { offset = 0 }

    if err := q.Order("id DESC").Limit(filter.Limit).Offset(offset).
        Preload("Photos").
        Preload("Toko").
        Preload("Category").
        Find(&items).Error; err != nil {
        return nil, 0, err
    }
    return items, count, nil
}

// Get product by ID with associations
func (r *Repository) GetByID(id uint) (*prodmodel.Product, error) {
    var p prodmodel.Product
    if err := r.db.Where("id = ?", id).
        Preload("Photos").Preload("Toko").Preload("Category").
        First(&p).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            return nil, nil
        }
        return nil, err
    }
    return &p, nil
}

func (r *Repository) GetBySlug(slug string) (*prodmodel.Product, error) {
    var p prodmodel.Product
    if err := r.db.Where("slug = ?", slug).First(&p).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &p, nil
}

// Create product and photos within transaction
func (r *Repository) Create(p *prodmodel.Product, photos []prodmodel.Photo) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Create(p).Error; err != nil {
            return err
        }
        if len(photos) > 0 {
            for i := range photos {
                photos[i].IDProduk = p.ID
            }
            if err := tx.Create(&photos).Error; err != nil {
                return err
            }
        }
        return nil
    })
}

// Update product fields (partial)
func (r *Repository) Update(p *prodmodel.Product, addPhotos []prodmodel.Photo) error {
    return r.db.Transaction(func(tx *gorm.DB) error {
        if err := tx.Model(&prodmodel.Product{}).Where("id = ?", p.ID).Updates(map[string]interface{}{
            "nama_produk":     p.NamaProduk,
            "slug":            p.Slug,
            "harga reseller":  p.HargaReseller,
            "harga konsumen":  p.HargaKonsumen,
            "stok":            p.Stok,
            "deskripsi":       p.Deskripsi,
            "id_toko":         p.IDToko,
            "id_category":     p.IDCategory,
        }).Error; err != nil {
            return err
        }
        if len(addPhotos) > 0 {
            for i := range addPhotos { addPhotos[i].IDProduk = p.ID }
            if err := tx.Create(&addPhotos).Error; err != nil { return err }
        }
        return nil
    })
}

func (r *Repository) Delete(id uint) error {
    return r.db.Delete(&prodmodel.Product{}, id).Error
}

// Ownership helpers
func (r *Repository) GetOwnerUserIDByProductID(id uint) (uint, error) {
    type row struct{ UserID uint }
    var out row
    // join produk -> toko to get toko.id_user
    err := r.db.Raw("SELECT t.id_user AS user_id FROM produk p JOIN toko t ON p.id_toko = t.id WHERE p.id = ?", id).Scan(&out).Error
    if err != nil { return 0, err }
    return out.UserID, nil
}

func (r *Repository) CategoryExists(id uint) (bool, error) {
    var cnt int64
    if err := r.db.Table("category").Where("id = ?", id).Count(&cnt).Error; err != nil {
        return false, err
    }
    return cnt > 0, nil
}

// Filter input for listing
type ListFilter struct {
    NamaProduk string
    CategoryID uint
    TokoID     uint
    MinHarga   *int
    MaxHarga   *int
    Limit      int
    Page       int
}