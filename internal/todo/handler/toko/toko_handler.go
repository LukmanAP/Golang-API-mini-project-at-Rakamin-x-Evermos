package toko

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	tokosvc "project-evermos/internal/todo/service/auth/toko"

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
	if u = strings.TrimSpace(strings.ToLower(u)); u == "" {
		return true // optional field
	}
	allowed := []string{".jpg", ".jpeg", ".png", ".gif", ".webp"}
	for _, ext := range allowed {
		if strings.HasSuffix(u, ext) {
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

	var payload struct {
		NamaToko string `json:"nama_toko"`
		UrlFoto  string `json:"url_foto"`
	}
	if err := c.BodyParser(&payload); err != nil {
		return fail(c, fiber.StatusBadRequest, "UPDATE", "Invalid JSON")
	}
	payload.NamaToko = strings.TrimSpace(payload.NamaToko)
	payload.UrlFoto = strings.TrimSpace(payload.UrlFoto)

	if len(payload.NamaToko) < 3 {
		return fail(c, fiber.StatusBadRequest, "UPDATE", "nama_toko minimal 3 karakter")
	}
	if !imgURLExtValid(payload.UrlFoto) {
		return fail(c, fiber.StatusBadRequest, "UPDATE", "url_foto harus URL file gambar (jpg|jpeg|png|gif|webp)")
	}

	if err := h.svc.UpdateStore(uint(id64), uid, payload.NamaToko, payload.UrlFoto); err != nil {
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

// JWTMiddleware validates JWT (HS256) from header: Authorization: Bearer <JWT>
// and sets Locals("user_id") for downstream handlers.
func JWTMiddleware(secret string) fiber.Handler {
    return func(c *fiber.Ctx) error {
        auth := strings.TrimSpace(c.Get("Authorization"))
        parts := strings.Fields(auth)
        if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") {
            return fail(c, fiber.StatusUnauthorized, "GET", "Unauthorized")
        }
        tok := strings.TrimSpace(parts[1])
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
