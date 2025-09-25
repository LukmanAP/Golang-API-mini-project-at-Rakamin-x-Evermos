package product

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"project-evermos/internal/config"
	prodmodel "project-evermos/internal/todo/model/product"
	tokoRepo "project-evermos/internal/todo/repository/toko"
	prodsvc "project-evermos/internal/todo/service/product"

	"github.com/gofiber/fiber/v2"
)

// Handler struct dan constructor
type Handler struct {
	s     *prodsvc.Service
	tokoR *tokoRepo.Repository
	cfg   *config.Config
}

func NewHandler(s *prodsvc.Service, tokoR *tokoRepo.Repository, cfg *config.Config) *Handler {
	return &Handler{s: s, tokoR: tokoR, cfg: cfg}
}

// Helpers untuk response standar
func respondFail(c *fiber.Ctx, httpStatus int, verb string, errs ...string) error {
	if len(errs) == 0 {
		errs = []string{"Unknown error"}
	}
	return c.Status(httpStatus).JSON(fiber.Map{
		"status":  false,
		"message": fmt.Sprintf("Failed to %s data", verb),
		"errors":  errs,
		"data":    nil,
	})
}

func respondOK(c *fiber.Ctx, verb string, data interface{}) error {
	return c.JSON(fiber.Map{
		"status":  true,
		"message": fmt.Sprintf("Succeed to %s data", verb),
		"errors":  nil,
		"data":    data,
	})
}

// Helper untuk ambil user_id dari JWT (c.Locals("user_id"))
func jwtUserID(c *fiber.Ctx) (uint, bool) {
	v := c.Locals("user_id")
	if v == nil {
		return 0, false
	}
	switch vv := v.(type) {
	case uint:
		return vv, true
	case int:
		if vv <= 0 {
			return 0, false
		}
		return uint(vv), true
	case int64:
		if vv <= 0 {
			return 0, false
		}
		return uint(vv), true
	case float64:
		if vv <= 0 {
			return 0, false
		}
		return uint(vv), true
	case string:
		if n, err := strconv.ParseUint(vv, 10, 64); err == nil {
			return uint(n), true
		}
	}
	return 0, false
}

// Endpoint: GET /product
func (h *Handler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	minStr := strings.TrimSpace(c.Query("min_harga", ""))
	maxStr := strings.TrimSpace(c.Query("max_harga", ""))
	var min, max *int
	if minStr != "" {
		if v, err := strconv.Atoi(minStr); err == nil {
			min = &v
		}
	}
	if maxStr != "" {
		if v, err := strconv.Atoi(maxStr); err == nil {
			max = &v
		}
	}

	items, _, err := h.s.List(prodsvc.ListParams{
		NamaProduk: c.Query("nama_produk", ""),
		CategoryID: parseUint(c.Query("category_id", "0")),
		TokoID:     parseUint(c.Query("toko_id", "0")),
		MinHarga:   min,
		MaxHarga:   max,
		Limit:      limit,
		Page:       page,
	})
	if err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "GET", err.Error())
	}

	out := make([]fiber.Map, 0, len(items))
	for _, p := range items {
		out = append(out, mapProductResponse(&p))
	}
	return respondOK(c, "GET", out)
}

// Endpoint: GET /product/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
	id := parseUint(c.Params("id"))
	p, err := h.s.GetByID(id)
	if err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "GET", err.Error())
	}
	if p == nil {
		return respondFail(c, fiber.StatusNotFound, "GET", "No Data Product")
	}
	return respondOK(c, "GET", mapProductResponse(p))
}

// Endpoint: POST /product (multipart)
func (h *Handler) Create(c *fiber.Ctx) error {
	uid, ok := jwtUserID(c)
	if !ok {
		return respondFail(c, fiber.StatusUnauthorized, "POST", "Unauthorized")
	}

	// validasi kepemilikan: user harus punya toko
	t, err := h.tokoR.FindByUserID(uid)
	if err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "POST", err.Error())
	}
	if t == nil {
		return respondFail(c, fiber.StatusBadRequest, "POST", "User belum memiliki toko")
	}

	name := strings.TrimSpace(c.FormValue("nama_produk"))
	catID := parseUint(c.FormValue("category_id"))
	hRes := atoiDefault(c.FormValue("harga_reseller"), -1)
	hKon := atoiDefault(c.FormValue("harga_konsumen"), -1)
	stok := atoiDefault(c.FormValue("stok"), -1)
	deskripsi := c.FormValue("deskripsi")

	// siapkan direktori upload
	if strings.TrimSpace(h.cfg.UploadDirProduct) != "" {
		_ = os.MkdirAll(h.cfg.UploadDirProduct, 0755)
	}

	// proses file
	form, ferr := c.MultipartForm()
	if ferr != nil && !strings.Contains(strings.ToLower(ferr.Error()), "multipart") {
		return respondFail(c, fiber.StatusBadRequest, "POST", "Invalid multipart form")
	}
	var savedURLs []string
	if form != nil {
		files := form.File["photos"]
		for _, f := range files {
			if !isAllowedImage(f.Filename) {
				return respondFail(c, fiber.StatusBadRequest, "POST", "Invalid image type")
			}
			if f.Size > 5*1024*1024 {
				return respondFail(c, fiber.StatusBadRequest, "POST", "File too large")
			}
			fname := fmt.Sprintf("%d-%s", time.Now().UnixNano(), sanitizeFilename(f.Filename))
			osPath := filepath.Join(h.cfg.UploadDirProduct, fname)
			urlPath := filepath.ToSlash(osPath)
			// pastikan direktori ada
			_ = os.MkdirAll(filepath.Dir(osPath), 0755)
			if err1 := c.SaveFile(f, osPath); err1 != nil {
				return respondFail(c, fiber.StatusInternalServerError, "POST", err.Error())
			}
			base := strings.TrimRight(h.cfg.BaseFileURL, "/")
			url := base + "/" + strings.TrimLeft(urlPath, "/")
			savedURLs = append(savedURLs, url)
		}
	}

	id, err := h.s.Create(prodsvc.CreateParams{
		UserID:        uid,
		NamaProduk:    name,
		CategoryID:    catID,
		HargaReseller: hRes,
		HargaKonsumen: hKon,
		Stok:          stok,
		Deskripsi:     deskripsi,
		PhotoURLs:     savedURLs,
		TokoID:        t.ID,
	})
	if err != nil {
		return respondFail(c, fiber.StatusBadRequest, "POST", err.Error())
	}
	return respondOK(c, "POST", id)
}

// Endpoint: PUT /product/:id (multipart update parsial)
func (h *Handler) Update(c *fiber.Ctx) error {
	uid, ok := jwtUserID(c)
	if !ok {
		return respondFail(c, fiber.StatusUnauthorized, "PUT", "Unauthorized")
	}

	id := parseUint(c.Params("id"))
	ownerID, err := h.s.RepoOwnerUserID(id)
	if err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "PUT", err.Error())
	}
	if ownerID == 0 {
		return respondFail(c, fiber.StatusNotFound, "PUT", "No Data Product")
	}
	if ownerID != uid {
		return respondFail(c, fiber.StatusForbidden, "PUT", "Tidak memiliki izin mengelola produk ini")
	}

	var namePtr *string
	if v := strings.TrimSpace(c.FormValue("nama_produk")); v != "" {
		namePtr = &v
	}
	var catPtr *uint
	if v := c.FormValue("category_id"); strings.TrimSpace(v) != "" {
		vv := parseUint(v)
		catPtr = &vv
	}
	var hResPtr *int
	if v := c.FormValue("harga_reseller"); strings.TrimSpace(v) != "" {
		n := atoiDefault(v, 0)
		hResPtr = &n
	}
	var hKonPtr *int
	if v := c.FormValue("harga_konsumen"); strings.TrimSpace(v) != "" {
		n := atoiDefault(v, 0)
		hKonPtr = &n
	}
	var stokPtr *int
	if v := c.FormValue("stok"); strings.TrimSpace(v) != "" {
		n := atoiDefault(v, 0)
		stokPtr = &n
	}
	var deskPtr *string
	if v := c.FormValue("deskripsi"); v != "" {
		deskPtr = &v
	}

	// siapkan direktori upload
	if strings.TrimSpace(h.cfg.UploadDirProduct) != "" {
		_ = os.MkdirAll(h.cfg.UploadDirProduct, 0755)
	}

	var savedURLs []string
	form, _ := c.MultipartForm()
	if form != nil {
		files := form.File["photos"]
		for _, f := range files {
			if !isAllowedImage(f.Filename) {
				return respondFail(c, fiber.StatusBadRequest, "PUT", "Invalid image type")
			}
			if f.Size > 5*1024*1024 {
				return respondFail(c, fiber.StatusBadRequest, "PUT", "File too large")
			}
			fname := fmt.Sprintf("%d-%s", time.Now().UnixNano(), sanitizeFilename(f.Filename))
			osPath := filepath.Join(h.cfg.UploadDirProduct, fname)
			urlPath := filepath.ToSlash(osPath)
			_ = os.MkdirAll(filepath.Dir(osPath), 0755)
			if err := c.SaveFile(f, osPath); err != nil {
				return respondFail(c, fiber.StatusInternalServerError, "PUT", err.Error())
			}
			base := strings.TrimRight(h.cfg.BaseFileURL, "/")
			url := base + "/" + strings.TrimLeft(urlPath, "/")
			savedURLs = append(savedURLs, url)
		}
	}

	if err := h.s.Update(prodsvc.UpdateParams{
		UserID:        uid,
		ID:            id,
		NamaProduk:    namePtr,
		CategoryID:    catPtr,
		HargaReseller: hResPtr,
		HargaKonsumen: hKonPtr,
		Stok:          stokPtr,
		Deskripsi:     deskPtr,
		PhotoURLs:     savedURLs,
	}); err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "not found") {
			return respondFail(c, fiber.StatusNotFound, "PUT", "No Data Product")
		}
		return respondFail(c, fiber.StatusBadRequest, "PUT", err.Error())
	}

	return respondOK(c, "PUT", "")
}

// Endpoint: DELETE /product/:id
func (h *Handler) Delete(c *fiber.Ctx) error {
	uid, ok := jwtUserID(c)
	if !ok {
		return respondFail(c, fiber.StatusUnauthorized, "DELETE", "Unauthorized")
	}

	id := parseUint(c.Params("id"))
	ownerID, err := h.s.RepoOwnerUserID(id)
	if err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "DELETE", err.Error())
	}
	if ownerID == 0 {
		return respondFail(c, fiber.StatusBadRequest, "DELETE", "record not found")
	}
	if ownerID != uid {
		return respondFail(c, fiber.StatusForbidden, "DELETE", "Tidak memiliki izin mengelola produk ini")
	}

	if err := h.s.Delete(id); err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "DELETE", err.Error())
	}
	return respondOK(c, "DELETE", "")
}

// Mapper respons produk
func mapProductResponse(p *prodmodel.Product) fiber.Map {
	photos := make([]fiber.Map, 0, len(p.Photos))
	for _, ph := range p.Photos {
		photos = append(photos, fiber.Map{"id": ph.ID, "product_id": ph.IDProduk, "url": ph.URL})
	}
	var toko fiber.Map
	if p.Toko != nil {
		toko = fiber.Map{"id": p.Toko.ID, "nama_toko": p.Toko.NamaToko, "url_foto": p.Toko.UrlFoto}
	} else {
		toko = fiber.Map{"id": p.IDToko}
	}
	var category fiber.Map
	if p.Category != nil {
		category = fiber.Map{"id": p.Category.ID, "nama_category": p.Category.NamaCategory}
	} else {
		category = fiber.Map{"id": p.IDCategory}
	}
	// harga di DB string -> ubah ke int untuk output
	hargaRes := atoiSafe(p.HargaReseller)
	hargaKon := atoiSafe(p.HargaKonsumen)

	return fiber.Map{
		"id":             p.ID,
		"nama_produk":    p.NamaProduk,
		"slug":           p.Slug,
		"harga_reseller": hargaRes,
		"harga_konsumen": hargaKon,
		"stok":           p.Stok,
		"deskripsi":      p.Deskripsi,
		"toko":           toko,
		"category":       category,
		"photos":         photos,
	}
}

// Utilities
func parseUint(s string) uint {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	if n < 0 {
		n = 0
	}
	return uint(n)
}
func atoiDefault(s string, def int) int {
	n, err := strconv.Atoi(strings.TrimSpace(s))
	if err != nil {
		return def
	}
	return n
}
func atoiSafe(s string) int {
	n, _ := strconv.Atoi(strings.TrimSpace(s))
	return n
}
func isAllowedImage(name string) bool {
	lower := strings.ToLower(name)
	return strings.HasSuffix(lower, ".jpg") || strings.HasSuffix(lower, ".jpeg") || strings.HasSuffix(lower, ".png")
}
func sanitizeFilename(name string) string {
	name = strings.ReplaceAll(name, " ", "-")
	name = strings.ReplaceAll(name, "..", "")
	name = strings.Trim(name, "-")
	return name
}
