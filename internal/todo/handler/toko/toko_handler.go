package toko

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	tokosvc "project-evermos/internal/todo/service/toko"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
)

// NewHandler constructs toko HTTP handlers
func NewHandler(s *tokosvc.Service) *Handler { return &Handler{svc: s} }

type Handler struct {
	svc *tokosvc.Service
}

// ---------- Helpers ----------

// jwtUserID extracts user ID set by middleware
func jwtUserID(c *fiber.Ctx) (uint, bool) {
	v := c.Locals("user_id")
	if v == nil {
		return 0, false
	}
	switch vv := v.(type) {
	case uint:
		return vv, true
	case int:
		if vv < 0 {
			return 0, false
		}
		return uint(vv), true
	case int64:
		if vv < 0 {
			return 0, false
		}
		return uint(vv), true
	case float64:
		if vv < 0 {
			return 0, false
		}
		return uint(vv), true
	case string:
		if vv == "" {
			return 0, false
		}
		if n, err := strconv.ParseUint(vv, 10, 64); err == nil {
			return uint(n), true
		}
	}
	return 0, false
}

func imgURLExtValid(u string) bool {
	u = strings.TrimSpace(strings.ToLower(u))
	if u == "" {
		return true // optional field
	}
	// ignore query params when checking extension
	if i := strings.Index(u, "?"); i >= 0 {
		u = u[:i]
	}
	allowed := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range allowed {
		if strings.HasSuffix(u, ext) {
			return true
		}
	}
	return false
}

func isAllowedImageFilename(name string) bool {
	name = strings.ToLower(strings.TrimSpace(name))
	allowed := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	ext := filepath.Ext(name)
	for _, e := range allowed {
		if ext == e {
			return true
		}
	}
	return false
}

func fail(c *fiber.Ctx, httpStatus int, verb string, errs ...string) error {
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

// ---------- Handlers ----------

// GET /toko/my
func (h *Handler) GetMy(c *fiber.Ctx) error {
	uid, ok := jwtUserID(c)
	if !ok {
		return fail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
	}
	data, err := h.svc.GetMyStore(uid)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "GET", err.Error())
	}
	// If user has no store, return data: null (as per policy)
	return respondOK(c, "GET", data)
}

// PUT /toko/:id_toko
func (h *Handler) Update(c *fiber.Ctx) error {
	uid, ok := jwtUserID(c)
	if !ok {
		return fail(c, fiber.StatusUnauthorized, "UPDATE", "Unauthorized")
	}
	idStr := c.Params("id_toko")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		return fail(c, fiber.StatusBadRequest, "UPDATE", "id_toko tidak valid")
	}

	// Support both JSON and form-data (multipart or x-www-form-urlencoded)
	var namaToko, urlFoto string
	ct := strings.ToLower(c.Get("Content-Type"))
	if strings.Contains(ct, "multipart/form-data") || strings.Contains(ct, "application/x-www-form-urlencoded") {
		// try file first
		if f, ferr := c.FormFile("photo"); ferr == nil && f != nil && f.Size > 0 {
			if !isAllowedImageFilename(f.Filename) {
				return fail(c, fiber.StatusBadRequest, "UPDATE", "photo harus file gambar (jpg|jpeg|png|gif|webp)")
			}
			_ = os.MkdirAll("uploads/stores", 0755)
			fname := fmt.Sprintf("toko-%d%s", time.Now().UnixNano(), filepath.Ext(f.Filename))
			dest := filepath.Join("uploads/stores", fname)
			if err := c.SaveFile(f, dest); err != nil {
				return fail(c, fiber.StatusInternalServerError, "UPDATE", "gagal menyimpan file foto")
			}
			// store relative path with forward slashes for consistency
			urlFoto = filepath.ToSlash(dest)
		} else {
			urlFoto = c.FormValue("photo")
		}
		namaToko = c.FormValue("nama_toko")
	} else {
		var payload struct {
			NamaToko string `json:"nama_toko"`
			Photo    string `json:"photo"`
		}
		if err := c.BodyParser(&payload); err != nil {
			return fail(c, fiber.StatusBadRequest, "UPDATE", "Invalid JSON or form data")
		}
		namaToko = payload.NamaToko
		urlFoto = payload.Photo
	}
	
	namaToko = strings.TrimSpace(namaToko)
	urlFoto = strings.TrimSpace(urlFoto)

	if len(namaToko) < 3 {
		return fail(c, fiber.StatusBadRequest, "UPDATE", "nama_toko minimal 3 karakter")
	}
	// If urlFoto provided as URL string, validate extension
	if urlFoto != "" && strings.HasPrefix(urlFoto, "http") && !imgURLExtValid(urlFoto) {
		return fail(c, fiber.StatusBadRequest, "UPDATE", "photo harus URL file gambar (jpg|jpeg|png|gif|webp)")
	}

	if err := h.svc.UpdateStore(uint(id64), uid, namaToko, urlFoto); err != nil {
		switch {
		case errors.Is(err, tokosvc.ErrNotFound):
			return fail(c, fiber.StatusNotFound, "UPDATE", "Toko tidak ditemukan")
		case errors.Is(err, tokosvc.ErrForbidden):
			return fail(c, fiber.StatusForbidden, "UPDATE", "Tidak memiliki izin mengelola toko ini")
		default:
			return fail(c, fiber.StatusInternalServerError, "UPDATE", err.Error())
		}
	}
	return respondOK(c, "UPDATE", "Update Succeed")
}

// GET /toko/:id_toko (public)
func (h *Handler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id_toko")
	id64, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil || id64 == 0 {
		return fail(c, fiber.StatusBadRequest, "GET", "id_toko tidak valid")
	}
	data, err := h.svc.GetByID(uint(id64), true, 0)
	if err != nil {
		switch {
		case errors.Is(err, tokosvc.ErrNotFound):
			return fail(c, fiber.StatusNotFound, "GET", "Toko tidak ditemukan")
		case errors.Is(err, tokosvc.ErrForbidden):
			return fail(c, fiber.StatusForbidden, "GET", "Tidak memiliki izin mengelola toko ini")
		default:
			return fail(c, fiber.StatusInternalServerError, "GET", err.Error())
		}
	}
	return respondOK(c, "GET", data)
}

// GET /toko?limit=&page=&nama=
func (h *Handler) List(c *fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "10"))
	page, _ := strconv.Atoi(c.Query("page", "1"))
	name := c.Query("nama", "")

	data, err := h.svc.List(limit, page, name)
	if err != nil {
		return fail(c, fiber.StatusInternalServerError, "GET", err.Error())
	}
	return respondOK(c, "GET", data)
}

// ---------- Middleware ----------

// JWTMiddleware validates JWT (HS256) from header. Supports:
// - token: <JWT>
// - Authorization: Bearer <JWT>
// and sets Locals("user_id") for downstream handlers.
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
            return fail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
        }
        tkn, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"]) 
            }
            return []byte(secret), nil
        })
        if err != nil || !tkn.Valid {
            return fail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
        }
        if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
            var uid uint
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
                return fail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
            }
            c.Locals("user_id", uid)
            return c.Next()
        }
        return fail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
    }
}
