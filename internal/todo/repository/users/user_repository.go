package users

import (
    "strings"

    usermodel "project-evermos/internal/todo/model/users"

    "gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

// ----- User queries -----
func (r *Repository) FindByID(id uint) (*usermodel.User, error) {
    var u usermodel.User
    if err := r.db.Where("id = ?", id).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *Repository) FindByEmail(email string) (*usermodel.User, error) {
    var u usermodel.User
    if err := r.db.Where("email = ?", strings.TrimSpace(email)).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *Repository) FindByPhone(notelp string) (*usermodel.User, error) {
    var u usermodel.User
    if err := r.db.Where("notelp = ?", strings.TrimSpace(notelp)).First(&u).Error; err != nil {
        return nil, err
    }
    return &u, nil
}

func (r *Repository) UpdateSelf(u *usermodel.User) error {
    return r.db.Model(&usermodel.User{}).Where("id = ?", u.ID).
        Updates(map[string]interface{}{
            "nama":           u.Nama,
            "kata_sandi":     u.KataSandi,
            "notelp":         u.NoTelp,
            "tanggal lahir":  u.TanggalLahir,
            "pekerjaan":      u.Pekerjaan,
            "email":          u.Email,
            "id_provinsi":    u.IDProvinsi,
            "id_kota":        u.IDKota,
        }).Error
}

// ----- Alamat queries (ownership enforced via where id_user = ?) -----
func (r *Repository) ListAlamatByUser(userID uint, titleFilter string) ([]usermodel.Alamat, error) {
    var rows []usermodel.Alamat
    q := r.db.Where("id_user = ?", userID)
    t := strings.TrimSpace(titleFilter)
    if t != "" {
        like := "%" + strings.ToLower(t) + "%"
        q = q.Where("LOWER(`judul alamat`) LIKE ?", like)
    }
    if err := q.Find(&rows).Error; err != nil { return nil, err }
    return rows, nil
}

func (r *Repository) GetAlamatByIDForUser(userID, alamatID uint) (*usermodel.Alamat, error) {
    var a usermodel.Alamat
    if err := r.db.Where("id_user = ? AND id = ?", userID, alamatID).First(&a).Error; err != nil {
        return nil, err
    }
    return &a, nil
}

// Fetch alamat by ID regardless of owner, for ownership checks
func (r *Repository) GetAlamatByID(id uint) (*usermodel.Alamat, error) {
    var a usermodel.Alamat
    if err := r.db.Where("id = ?", id).First(&a).Error; err != nil {
        return nil, err
    }
    return &a, nil
}

func (r *Repository) CreateAlamat(a *usermodel.Alamat) error {
    return r.db.Create(a).Error
}

func (r *Repository) UpdateAlamatForUser(userID uint, a *usermodel.Alamat) error {
    return r.db.Model(&usermodel.Alamat{}).
        Where("id_user = ? AND id = ?", userID, a.ID).
        Updates(map[string]interface{}{
            "nama penerima": a.NamaPenerima,
            "no telp":      a.NoTelp,
            "detail_alamat": a.DetailAlamat,
        }).Error
}

func (r *Repository) DeleteAlamatForUser(userID, alamatID uint) error {
    return r.db.Where("id_user = ? AND id = ?", userID, alamatID).Delete(&usermodel.Alamat{}).Error
}