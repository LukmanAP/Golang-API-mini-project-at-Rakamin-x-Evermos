package users

import (
	"errors"
	"fmt"
	"net/mail"
	"strconv"
	"strings"
	"time"

	svc "project-evermos/internal/todo/service/users"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
)

type Handler struct{ s *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{s: s} }

// --- helpers ---
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

// --- middleware ---
// Per requirement: baca header 'token: <JWT>' (bukan Authorization: Bearer)
func JWTMiddleware(secret string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Try custom header 'token' first
		tok := strings.TrimSpace(c.Get("token"))
		// Fallback: support standard Authorization: Bearer <token>
		if tok == "" {
			auth := strings.TrimSpace(c.Get("Authorization"))
			if auth != "" {
				parts := strings.Fields(auth)
				if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
					tok = strings.TrimSpace(parts[1])
				}
			}
		}
		if tok == "" {
			return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
		}
		tkn, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]) 
			}
			return []byte(secret), nil
		})
		if err != nil || !tkn.Valid {
			return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
		}
		if claims, okc := tkn.Claims.(jwt.MapClaims); okc {
			var uid uint
			// coba beberapa key umum pada claims
			if v, ok := claims["user_id"]; ok {
				switch vv := v.(type) {
				case float64:
					uid = uint(vv)
				case string:
					if n, err := strconv.ParseUint(vv, 10, 64); err == nil {
						uid = uint(n)
					}
				}
			} else if v, ok := claims["id"]; ok {
				switch vv := v.(type) {
				case float64:
					uid = uint(vv)
				case string:
					if n, err := strconv.ParseUint(vv, 10, 64); err == nil {
						uid = uint(n)
					}
				}
			} else if v, ok := claims["sub"]; ok {
				if s, ok := v.(string); ok {
					if n, err := strconv.ParseUint(s, 10, 64); err == nil {
						uid = uint(n)
					}
				}
			}
			if uid == 0 {
				return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
			}
			c.Locals("user_id", uid)
			return c.Next()
		}
		return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
}

// --- handlers ---
// GET /user
func (h *Handler) GetProfile(c *fiber.Ctx) error {
	uid, okJWT := jwtUserID(c)
	if !okJWT {
		return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
	u, err := h.s.GetProfile(uid)
	if err != nil {
		return respondFail(c, fiber.StatusNotFound, "GET", "record not found")
	}
	// Format tanggal lahir dd/MM/yyyy
	var tgl string
	if u.TanggalLahir != nil {
		tgl = u.TanggalLahir.Format("02/01/2006")
	} else {
		tgl = ""
	}
	return respondOK(c, "GET", fiber.Map{
		"id":            u.ID,
		"nama":          u.Nama,
		"no_telp":       u.NoTelp,
		"tanggal_Lahir": tgl,
		"pekerjaan":     u.Pekerjaan,
		"email":         u.Email,
		"id_provinsi":   fiber.Map{"id": u.IDProvinsi, "name": provinceName(u.IDProvinsi)},
		"id_kota":       fiber.Map{"id": u.IDKota, "province_id": u.IDProvinsi, "name": cityName(u.IDKota)},
	})
}

// PUT /user
func (h *Handler) UpdateProfile(c *fiber.Ctx) error {
	uid, okJWT := jwtUserID(c)
	if !okJWT {
		return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
	var body struct {
		Nama         string `json:"nama"`
		KataSandi    string `json:"kata_sandi"`
		NoTelp       string `json:"no_telp"`
		TanggalLahir string `json:"tanggal_Lahir"`
		Pekerjaan    string `json:"pekerjaan"`
		Email        string `json:"email"`
		IDProvinsi   string `json:"id_provinsi"`
		IDKota       string `json:"id_kota"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondFail(c, fiber.StatusBadRequest, "GET", "Invalid JSON")
	}
	if strings.TrimSpace(body.Email) != "" {
		if _, err := mail.ParseAddress(strings.TrimSpace(body.Email)); err != nil {
			return respondFail(c, fiber.StatusBadRequest, "GET", "Email tidak valid")
		}
	}
	if strings.TrimSpace(body.TanggalLahir) != "" {
		// validasi format dd/MM/yyyy
		if _, err := time.Parse("02/01/2006", body.TanggalLahir); err != nil {
			return respondFail(c, fiber.StatusBadRequest, "GET", "tanggal_Lahir format harus dd/MM/yyyy")
		}
	}
	if err := h.s.UpdateProfile(uid, svc.UpdateProfileInput{
		Nama:         body.Nama,
		KataSandi:    body.KataSandi,
		NoTelp:       body.NoTelp,
		TanggalLahir: body.TanggalLahir,
		Pekerjaan:    body.Pekerjaan,
		Email:        body.Email,
		IDProvinsi:   body.IDProvinsi,
		IDKota:       body.IDKota,
	}); err != nil {
		msg := strings.ToLower(err.Error())
		if strings.Contains(msg, "duplicate") || strings.Contains(msg, "unique") || strings.Contains(msg, "1062") {
			if strings.Contains(msg, "email") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  false,
					"message": "Failed to GET data",
					"errors":  []string{"Error 1062: Duplicate entry '" + strings.TrimSpace(body.Email) + "' for key 'users.email'"},
					"data":    nil,
				})
			}
			if strings.Contains(msg, "notelp") || strings.Contains(msg, "phone") || strings.Contains(msg, "no_telp") {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"status":  false,
					"message": "Failed to GET data",
					"errors":  []string{"Error 1062: Duplicate entry '" + strings.TrimSpace(body.NoTelp) + "' for key 'users.notelp'"},
					"data":    nil,
				})
			}
			return respondFail(c, fiber.StatusBadRequest, "GET", "duplicate")
		}
		if strings.Contains(msg, "tanggal_lahir") {
			return respondFail(c, fiber.StatusBadRequest, "GET", "tanggal_Lahir format harus dd/MM/yyyy")
		}
		return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
	}
	return respondOK(c, "GET", "")
}

// GET /user/alamat?judul_alamat=
func (h *Handler) ListAlamat(c *fiber.Ctx) error {
	uid, okJWT := jwtUserID(c)
	if !okJWT {
		return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
	title := c.Query("judul_alamat", "")
	rows, err := h.s.ListAlamat(uid, title)
	if err != nil {
		return respondFail(c, fiber.StatusInternalServerError, "GET", err.Error())
	}
	out := make([]fiber.Map, 0, len(rows))
	for _, a := range rows {
		out = append(out, fiber.Map{
			"id":            a.ID,
			"judul_alamat":  a.JudulAlamat,
			"nama_penerima": a.NamaPenerima,
			"no_telp":       a.NoTelp,
			"detail_alamat": a.DetailAlamat,
		})
	}
	return respondOK(c, "GET", out)
}

// GET /user/alamat/:id
func (h *Handler) GetAlamat(c *fiber.Ctx) error {
	uid, okJWT := jwtUserID(c)
	if !okJWT {
		return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		return respondFail(c, fiber.StatusBadRequest, "GET", "id tidak valid")
	}
	a, err := h.s.GetAlamat(uid, uint(id64))
	if err != nil {
		if errors.Is(err, svc.ErrNotFound) {
			return respondFail(c, fiber.StatusNotFound, "GET", "record not found")
		}
		if errors.Is(err, svc.ErrForbidden) {
			return respondFail(c, fiber.StatusForbidden, "GET", "forbidden")
		}
		return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
	}
	return respondOK(c, "GET", fiber.Map{
		"id":            a.ID,
		"judul_alamat":  a.JudulAlamat,
		"nama_penerima": a.NamaPenerima,
		"no_telp":       a.NoTelp,
		"detail_alamat": a.DetailAlamat,
	})
}

// POST /user/alamat
func (h *Handler) CreateAlamat(c *fiber.Ctx) error {
	uid, okJWT := jwtUserID(c)
	if !okJWT {
		return respondFail(c, fiber.StatusUnauthorized, "POST", "Unauthorized")
	}
	var body struct {
		JudulAlamat  string `json:"judul_alamat"`
		NamaPenerima string `json:"nama_penerima"`
		NoTelp       string `json:"no_telp"`
		DetailAlamat string `json:"detail_alamat"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondFail(c, fiber.StatusBadRequest, "POST", "Invalid JSON")
	}
	id, err := h.s.CreateAlamat(uid, svc.CreateAlamatInput{
		JudulAlamat:  body.JudulAlamat,
		NamaPenerima: body.NamaPenerima,
		NoTelp:       body.NoTelp,
		DetailAlamat: body.DetailAlamat,
	})
	if err != nil {
		return respondFail(c, fiber.StatusBadRequest, "POST", err.Error())
	}
	return respondOK(c, "POST", id)
}

// PUT /user/alamat/:id
func (h *Handler) UpdateAlamat(c *fiber.Ctx) error {
    uid, okJWT := jwtUserID(c)
    if !okJWT {
        return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
    }
    idStr := c.Params("id")
    id64, err := strconv.ParseUint(idStr, 10, 64)
    if err != nil || id64 == 0 {
        return respondFail(c, fiber.StatusBadRequest, "GET", "id tidak valid")
    }
    var body struct {
        JudulAlamat  string `json:"judul_alamat"`
        NamaPenerima string `json:"nama_penerima"`
        NoTelp       string `json:"no_telp"`
        DetailAlamat string `json:"detail_alamat"`
    }
    if err := c.BodyParser(&body); err != nil {
        return respondFail(c, fiber.StatusBadRequest, "GET", "Invalid JSON")
    }
    if err := h.s.UpdateAlamat(uid, uint(id64), svc.UpdateAlamatInput{
        JudulAlamat:  body.JudulAlamat,
        NamaPenerima: body.NamaPenerima,
        NoTelp:       body.NoTelp,
        DetailAlamat: body.DetailAlamat,
    }); err != nil {
        if errors.Is(err, svc.ErrNotFound) {
            return respondFail(c, fiber.StatusNotFound, "GET", "record not found")
        }
        if errors.Is(err, svc.ErrForbidden) {
            return respondFail(c, fiber.StatusForbidden, "GET", "forbidden")
        }
        return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
    }
    return respondOK(c, "GET", "")
}

// DELETE /user/alamat/:id
func (h *Handler) DeleteAlamat(c *fiber.Ctx) error {
	uid, okJWT := jwtUserID(c)
	if !okJWT {
		return respondFail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		return respondFail(c, fiber.StatusBadRequest, "GET", "id tidak valid")
	}
	if err := h.s.DeleteAlamat(uid, uint(id64)); err != nil {
		if errors.Is(err, svc.ErrNotFound) {
			return respondFail(c, fiber.StatusNotFound, "GET", "record not found")
		}
		if errors.Is(err, svc.ErrForbidden) {
			return respondFail(c, fiber.StatusForbidden, "GET", "forbidden")
		}
		return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
	}
	return respondOK(c, "GET", "")
}

// placeholder names; replace dengan lookup sebenarnya jika tersedia
func provinceName(id string) string { return id }
func cityName(id string) string     { return id }
