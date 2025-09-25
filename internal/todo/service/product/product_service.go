package product

import (
    "errors"
    "fmt"
    "path/filepath"
    "regexp"
    "strings"
    "time"

    prodmodel "project-evermos/internal/todo/model/product"
    prodrepo "project-evermos/internal/todo/repository/product"

    "gorm.io/gorm"
)

type Service struct {
    repo *prodrepo.Repository
    // baseURL used to construct file URLs for photos
    baseURL string
}

func NewService(repo *prodrepo.Repository, baseURL string) *Service {
    return &Service{repo: repo, baseURL: strings.TrimRight(baseURL, "/")}
}

type ListParams struct {
    NamaProduk string
    CategoryID uint
    TokoID     uint
    MinHarga   *int
    MaxHarga   *int
    Limit      int
    Page       int
}

type CreateParams struct {
    UserID        uint
    NamaProduk    string
    CategoryID    uint
    HargaReseller int
    HargaKonsumen int
    Stok          int
    Deskripsi     string
    // Photos holds already-saved file relative URLs (to be persisted)
    PhotoURLs []string
    // TokoID must be the user's toko id
    TokoID uint
}

type UpdateParams struct {
    UserID        uint
    ID            uint
    NamaProduk    *string
    CategoryID    *uint
    HargaReseller *int
    HargaKonsumen *int
    Stok          *int
    Deskripsi     *string
    PhotoURLs     []string // new photos to add
}

var (
    reNonWord        = regexp.MustCompile(`[^a-z0-9]+`)
    gormErrNotFound  = gorm.ErrRecordNotFound
)

func (s *Service) List(p ListParams) ([]prodmodel.Product, int64, error) {
    limit := p.Limit
    if limit <= 0 { limit = 10 }
    if limit > 100 { limit = 100 }
    page := p.Page
    if page <= 0 { page = 1 }
    f := prodrepo.ListFilter{
        NamaProduk: p.NamaProduk,
        CategoryID: p.CategoryID,
        TokoID:     p.TokoID,
        MinHarga:   p.MinHarga,
        MaxHarga:   p.MaxHarga,
        Limit:      limit,
        Page:       page,
    }
    return s.repo.List(f)
}

func (s *Service) GetByID(id uint) (*prodmodel.Product, error) {
    return s.repo.GetByID(id)
}

func (s *Service) Create(p CreateParams) (uint, error) {
    if err := s.validateCreate(p); err != nil { return 0, err }
    // build model
    prod := prodmodel.Product{
        NamaProduk:    p.NamaProduk,
        Slug:          "", // set below
        HargaReseller: fmt.Sprintf("%d", p.HargaReseller),
        HargaKonsumen: fmt.Sprintf("%d", p.HargaKonsumen),
        Stok:          p.Stok,
        Deskripsi:     p.Deskripsi,
        IDToko:        p.TokoID,
        IDCategory:    p.CategoryID,
        CreatedAt:     time.Now(),
        UpdatedAt:     time.Now(),
    }
    // generate unique slug
    slug, err := s.generateUniqueSlug(p.NamaProduk)
    if err != nil { return 0, err }
    prod.Slug = slug

    // convert urls to photo models
    var photos []prodmodel.Photo
    for _, url := range p.PhotoURLs {
        photos = append(photos, prodmodel.Photo{URL: url})
    }

    if err := s.repo.Create(&prod, photos); err != nil { return 0, err }
    return prod.ID, nil
}

func (s *Service) Update(p UpdateParams) error {
    // read existing product
    existing, err := s.repo.GetByID(p.ID)
    if err != nil { return err }
    if existing == nil { return gormErrNotFound }

    // apply partial updates
    if p.NamaProduk != nil {
        existing.NamaProduk = *p.NamaProduk
        slug, err := s.generateUniqueSlug(*p.NamaProduk)
        if err != nil { return err }
        existing.Slug = slug
    }
    if p.CategoryID != nil { existing.IDCategory = *p.CategoryID }
    if p.HargaReseller != nil { existing.HargaReseller = fmt.Sprintf("%d", *p.HargaReseller) }
    if p.HargaKonsumen != nil { existing.HargaKonsumen = fmt.Sprintf("%d", *p.HargaKonsumen) }
    if p.Stok != nil { existing.Stok = *p.Stok }
    if p.Deskripsi != nil { existing.Deskripsi = *p.Deskripsi }
    existing.UpdatedAt = time.Now()

    // photos to add
    var photos []prodmodel.Photo
    for _, url := range p.PhotoURLs {
        photos = append(photos, prodmodel.Photo{URL: url})
    }

    return s.repo.Update(existing, photos)
}

func (s *Service) Delete(id uint) error {
    return s.repo.Delete(id)
}

// Ownership helper for handlers
func (s *Service) RepoOwnerUserID(productID uint) (uint, error) {
    return s.repo.GetOwnerUserIDByProductID(productID)
}

// generateUniqueSlug builds slug from name and appends numeric suffix when clashing
func (s *Service) generateUniqueSlug(name string) (string, error) {
    base := strings.ToLower(strings.TrimSpace(name))
    base = reNonWord.ReplaceAllString(base, "-")
    base = strings.Trim(base, "-")
    if base == "" { base = fmt.Sprintf("produk-%d", time.Now().Unix()) }

    slug := base
    // naive loop: try up to reasonable attempts by checking in DB
    for i := 0; i < 1000; i++ {
        exists, err := s.slugExists(slug)
        if err != nil { return "", err }
        if !exists { return slug, nil }
        slug = fmt.Sprintf("%s-%d", base, i+1)
    }
    return "", errors.New("failed to generate unique slug")
}

func (s *Service) slugExists(slug string) (bool, error) {
    // check by query
    itm, err := s.repo.GetBySlug(slug)
    if err != nil { return false, err }
    return itm != nil, nil
}

// Helpers
func buildPhotoURL(baseURL, storedName string) string {
    if strings.HasPrefix(storedName, "http://") || strings.HasPrefix(storedName, "https://") {
        return storedName
    }
    return strings.TrimRight(baseURL, "/") + "/" + strings.TrimLeft(filepath.ToSlash(storedName), "/")
}

// Validation
func (s *Service) validateCreate(p CreateParams) error {
    var errs []string
    if len(strings.TrimSpace(p.NamaProduk)) < 3 {
        errs = append(errs, "nama_produk min 3 char")
    }
    if p.CategoryID == 0 {
        errs = append(errs, "category_id required")
    } else {
        ok, err := s.repo.CategoryExists(p.CategoryID)
        if err != nil { return err }
        if !ok { errs = append(errs, "category_id invalid") }
    }
    if p.HargaReseller < 0 || p.HargaKonsumen < 0 || p.Stok < 0 {
        errs = append(errs, "harga_reseller/harga_konsumen/stok must be >= 0")
    }
    if len(errs) > 0 { return errors.New(strings.Join(errs, "; ")) }
    return nil
}