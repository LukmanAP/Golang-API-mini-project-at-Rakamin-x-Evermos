package auth

import (
	"errors"
	"net/mail"
	"regexp"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	model "project-evermos/internal/todo/model/auth"
	service "project-evermos/internal/todo/service/auth"
)

// Handler handles auth-related HTTP requests.
type Handler struct{
	svc *service.Service
}

// NewHandler constructs a new Handler.
func NewHandler(s *service.Service) *Handler { return &Handler{svc: s} }

var (
	phoneRegex = regexp.MustCompile(`^\d{10,15}$`)
)

type loginRequest struct {
	NoTelp    string `json:"no_telp"`
	KataSandi string `json:"kata_sandi"`
}

type registerRequest struct {
	Nama         string `json:"nama"`
	KataSandi    string `json:"kata_sandi"`
	NoTelp       string `json:"no_telp"`
	TanggalLahir string `json:"tanggal_Lahir"`
	Pekerjaan    string `json:"pekerjaan"`
	Email        string `json:"email"`
	IDProvinsi   string `json:"id_provinsi"`
	IDKota       string `json:"id_kota"`
}

func (h *Handler) Login(c *fiber.Ctx) error {
	var req loginRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Invalid JSON"},
			"data":    nil,
		})
	}
	// validation
	errs := make([]string, 0)
	if !phoneRegex.MatchString(strings.TrimSpace(req.NoTelp)) {
		errs = append(errs, "No Telp tidak valid")
	}
	if len(req.KataSandi) < 6 {
		errs = append(errs, "Kata sandi minimal 6 karakter")
	}
	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  errs,
			"data":    nil,
		})
	}

	u, token, err := h.svc.Login(req.NoTelp, req.KataSandi)
	if err != nil {
		if errors.Is(err, service.ErrInvalidCredentials) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to POST data",
				"errors":  []string{"No Telp atau kata sandi salah"},
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	var tanggalLahir string
	if u.TanggalLahir != nil {
		tanggalLahir = u.TanggalLahir.Format("02/01/2006")
	} else {
		tanggalLahir = ""
	}

	resp := fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data": fiber.Map{
			"nama":           u.Nama,
			"no_telp":        u.NoTelp,
			"tanggal_Lahir":  tanggalLahir,
			"tentang":        valueOrEmpty(u.Tentang),
			"pekerjaan":      u.Pekerjaan,
			"email":          u.Email,
			"id_provinsi":    fiber.Map{"id": u.IDProvinsi, "name": provinceName(u.IDProvinsi)},
			"id_kota":        fiber.Map{"id": u.IDKota, "province_id": u.IDProvinsi, "name": cityName(u.IDKota)},
			"token":          token,
		},
	}
	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *Handler) Register(c *fiber.Ctx) error {
	var req registerRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{"Invalid JSON"},
			"data":    nil,
		})
	}

	// validation
	errs := make([]string, 0)
	if s := strings.TrimSpace(req.Nama); len(s) < 2 || len(s) > 100 {
		errs = append(errs, "Nama wajib 2–100 karakter")
	}
	if len(req.KataSandi) < 6 {
		errs = append(errs, "Kata sandi minimal 6 karakter")
	}
	if !phoneRegex.MatchString(strings.TrimSpace(req.NoTelp)) {
		errs = append(errs, "No Telp tidak valid")
	}
	if _, err := mail.ParseAddress(strings.TrimSpace(req.Email)); err != nil {
		errs = append(errs, "Email tidak valid")
	}
	if strings.TrimSpace(req.IDProvinsi) == "" || strings.TrimSpace(req.IDKota) == "" {
		errs = append(errs, "id_provinsi dan id_kota wajib diisi")
	}
	// verify city belongs to province (placeholder since we don't have tables); ensure IDs share prefix or pass basic check
	if !cityBelongsToProvince(req.IDKota, req.IDProvinsi) {
		errs = append(errs, "id_kota tidak sesuai dengan id_provinsi")
	}
	// parse tanggal lahir
	var tgl *time.Time
	if strings.TrimSpace(req.TanggalLahir) != "" {
		t, err := time.Parse("02/01/2006", req.TanggalLahir)
		if err != nil {
			errs = append(errs, "tanggal_Lahir tidak valid, gunakan format dd/MM/yyyy")
		} else {
			tgl = &t
		}
	}

	if s := strings.TrimSpace(req.Pekerjaan); len(s) < 2 || len(s) > 100 {
		errs = append(errs, "Pekerjaan wajib 2–100 karakter")
	}

	if len(errs) > 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  errs,
			"data":    nil,
		})
	}

	u := &model.User{
		Nama:         strings.TrimSpace(req.Nama),
		NoTelp:       strings.TrimSpace(req.NoTelp),
		TanggalLahir: tgl,
		Pekerjaan:    strings.TrimSpace(req.Pekerjaan),
		Email:        strings.TrimSpace(req.Email),
		IDProvinsi:   strings.TrimSpace(req.IDProvinsi),
		IDKota:       strings.TrimSpace(req.IDKota),
	}

	if err := h.svc.Register(u, req.KataSandi); err != nil {
		if errors.Is(err, service.ErrDuplicateEmail) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to POST data",
				"errors":  []string{"Error 1062: Duplicate entry '" + u.Email + "' for key 'users.email'"},
				"data":    nil,
			})
		}
		if errors.Is(err, service.ErrDuplicatePhone) {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"status":  false,
				"message": "Failed to POST data",
				"errors":  []string{"Error 1062: Duplicate entry '" + u.NoTelp + "' for key 'users.notelp'"},
				"data":    nil,
			})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"status":  false,
			"message": "Failed to POST data",
			"errors":  []string{err.Error()},
			"data":    nil,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to POST data",
		"errors":  nil,
		"data":    "Register Succeed",
	})
}

func valueOrEmpty(p *string) string {
	if p == nil { return "" }
	return *p
}

// cityBelongsToProvince verifies id_kota belongs to id_provinsi.
// Placeholder rule: assume kota starts with province id or both non-empty (since no reference table).
func cityBelongsToProvince(kota, prov string) bool {
	kota = strings.TrimSpace(kota)
	prov = strings.TrimSpace(prov)
	if kota == "" || prov == "" { return false }
	return strings.HasPrefix(kota, prov) || kota == prov
}

// provinceName and cityName are placeholders; in a real app, fetch from reference tables/APIs.
func provinceName(id string) string { return id }
func cityName(id string) string { return id }