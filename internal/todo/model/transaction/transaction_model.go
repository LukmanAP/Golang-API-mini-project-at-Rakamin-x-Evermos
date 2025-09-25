package transaction

import "time"

// Trx maps to trx table
type Trx struct {
    ID               uint       `gorm:"primaryKey;column:id"`
    IDUser           uint       `gorm:"column:id_user"`
    AlamatPengiriman uint       `gorm:"column:alamat_pengiriman"`
    HargaTotal       int        `gorm:"column:harga_total"`
    KodeInvoice      string     `gorm:"column:kode_invoice"`
    MethodBayar      string     `gorm:"column:method_bayar"`
    UpdatedAt        *time.Time `gorm:"column:updated_at"`
    CreatedAt        *time.Time `gorm:"column:created_at"`
}

func (Trx) TableName() string { return "trx" }

// DetailTrx maps to detail_trx table
type DetailTrx struct {
    ID          uint       `gorm:"primaryKey;column:id"`
    IDTrx       uint       `gorm:"column:id_trx"`
    IDLogProduk uint       `gorm:"column:id_log_produk"`
    IDToko      uint       `gorm:"column:id_toko"`
    Kuantitas   int        `gorm:"column:kuantitas"`
    HargaTotal  int        `gorm:"column:harga_total"`
    UpdatedAt   *time.Time `gorm:"column:updated_at"`
    CreatedAt   *time.Time `gorm:"column:created_at"`
}

func (DetailTrx) TableName() string { return "detail_trx" }

// LogProduk stores product snapshot at time of transaction
// Note: columns with spaces must be mapped carefully
type LogProduk struct {
    ID            uint       `gorm:"primaryKey;column:id"`
    IDProduk      uint       `gorm:"column:id_produk"`
    NamaProduk    string     `gorm:"column:nama_produk"`
    Slug          string     `gorm:"column:slug"`
    HargaReseller string     `gorm:"column:harga reseller"`
    HargaKonsumen string     `gorm:"column:harga konsumen"`
    Deskripsi     string     `gorm:"column:deskripsi"`
    IDToko        uint       `gorm:"column:id_toko"`
    IDCategory    uint       `gorm:"column:id_category"`
    PhotosJSON    string     `gorm:"column:photos_json"` // JSON array of photo URLs
    UpdatedAt     *time.Time `gorm:"column:updated_at"`
    CreatedAt     *time.Time `gorm:"column:created_at"`
}

func (LogProduk) TableName() string { return "log_produk" }