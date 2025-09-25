package transaction

import (
    "errors"
    "encoding/json"

    prodmodel "project-evermos/internal/todo/model/product"
    tokomodel "project-evermos/internal/todo/model/toko"
    usermodel "project-evermos/internal/todo/model/users"
    trxmodel "project-evermos/internal/todo/model/transaction"

    "gorm.io/gorm"
)

// Repository provides data access for transaction domain.
type Repository struct { DB *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{DB: db} }

// --- create helpers within a transaction ---
func (r *Repository) CreateTrx(tx *gorm.DB, t *trxmodel.Trx) error {
    return tx.Create(t).Error
}

func (r *Repository) CreateDetailItems(tx *gorm.DB, items []trxmodel.DetailTrx) error {
    if len(items) == 0 { return nil }
    return tx.Create(&items).Error
}

func (r *Repository) CreateLogProduk(tx *gorm.DB, lp *trxmodel.LogProduk) error {
    return tx.Create(lp).Error
}

// UpdateProductStock decrements stock safely
func (r *Repository) UpdateProductStock(tx *gorm.DB, productID uint, dec int) error {
    if dec <= 0 { return nil }
    res := tx.Exec("UPDATE produk SET stok = stok - ? WHERE id = ? AND stok >= ?", dec, productID, dec)
    if res.Error != nil { return res.Error }
    if res.RowsAffected != 1 { return errors.New("insufficient stock") }
    return nil
}

// --- fetch helpers ---
func (r *Repository) GetProductByID(id uint) (*prodmodel.Product, error) {
    var p prodmodel.Product
    if err := r.DB.Where("id = ?", id).Preload("Photos").First(&p).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &p, nil
}

func (r *Repository) GetAlamatByID(id uint) (*usermodel.Alamat, error) {
    var a usermodel.Alamat
    if err := r.DB.Where("id = ?", id).First(&a).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &a, nil
}

func (r *Repository) GetOwnerUserIDOfTrx(trxID uint) (uint, error) {
    type row struct{ UID uint }
    var out row
    err := r.DB.Raw("SELECT id_user AS uid FROM trx WHERE id = ?", trxID).Scan(&out).Error
    if err != nil { return 0, err }
    return out.UID, nil
}

func (r *Repository) GetTrxByID(id uint) (*trxmodel.Trx, error) {
    var t trxmodel.Trx
    if err := r.DB.Where("id = ?", id).First(&t).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &t, nil
}

func (r *Repository) ListTrxByUser(userID uint, limit, page int) ([]trxmodel.Trx, int64, error) {
    var rows []trxmodel.Trx
    var cnt int64
    q := r.DB.Model(&trxmodel.Trx{}).Where("id_user = ?", userID)
    if err := q.Count(&cnt).Error; err != nil { return nil, 0, err }
    if limit <= 0 { limit = 10 }
    if limit > 100 { limit = 100 }
    if page <= 0 { page = 1 }
    off := (page - 1) * limit
    if err := q.Order("id DESC").Limit(limit).Offset(off).Find(&rows).Error; err != nil {
        return nil, 0, err
    }
    return rows, cnt, nil
}

func (r *Repository) GetDetailItems(trxID uint) ([]trxmodel.DetailTrx, error) {
    var items []trxmodel.DetailTrx
    if err := r.DB.Where("id_trx = ?", trxID).Find(&items).Error; err != nil { return nil, err }
    return items, nil
}

func (r *Repository) GetLogProdukByID(id uint) (*trxmodel.LogProduk, error) {
    var lp trxmodel.LogProduk
    if err := r.DB.Where("id = ?", id).First(&lp).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &lp, nil
}

func (r *Repository) GetTokoByID(id uint) (*tokomodel.Toko, error) {
    var t tokomodel.Toko
    if err := r.DB.Where("id = ?", id).First(&t).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &t, nil
}

func (r *Repository) GetCategoryByID(id uint) (*prodmodel.CategoryRef, error) {
    var c prodmodel.CategoryRef
    if err := r.DB.Where("id = ?", id).First(&c).Error; err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) { return nil, nil }
        return nil, err
    }
    return &c, nil
}

// Helpers
func MarshalPhotos(urls []string) string {
    b, _ := json.Marshal(urls)
    return string(b)
}