package category

import (
	"errors"
	"strconv"
	"strings"

	usersRepo "project-evermos/internal/todo/repository/users"
	svc "project-evermos/internal/todo/service/category"

	"github.com/gofiber/fiber/v2"
	jwt "github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

type Handler struct{ s *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{s: s} }

// -------- Helpers for response contract --------
func respondOK(c *fiber.Ctx, op string, data interface{}) error {
	msg := "Succeed to GET data"
	if strings.EqualFold(op, "POST") {
		msg = "Succeed to POST data"
	}
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": msg,
		"errors":  nil,
		"data":    data,
	})
}

func respondFail(c *fiber.Ctx, code int, op string, errs ...string) error {
	msg := "Failed to GET data"
	if strings.EqualFold(op, "POST") {
		msg = "Failed to POST data"
	}
	if len(errs) == 0 {
		errs = []string{"Unknown error"}
	}
	return c.Status(code).JSON(fiber.Map{
		"status":  false,
		"message": msg,
		"errors":  errs,
		"data":    nil,
	})
}

// -------- Middleware --------
// AuthJWT reads header 'token' (or Authorization: Bearer) and sets user_id & is_admin in Context.
// On invalid/absent token, returns 401 with Failed to POST data.
func AuthJWT(secret string, db *gorm.DB, op string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		tok := strings.TrimSpace(c.Get("token"))
		if tok == "" {
			// fallback Bearer
			auth := strings.TrimSpace(c.Get("Authorization"))
			parts := strings.Fields(auth)
			if len(parts) == 2 && strings.EqualFold(parts[0], "Bearer") {
				tok = strings.TrimSpace(parts[1])
			}
		}
		if tok == "" {
			return respondFail(c, fiber.StatusUnauthorized, op, "Unauthorized")
		}
		tkn, err := jwt.Parse(tok, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return []byte(secret), nil
		})
		if err != nil || !tkn.Valid {
			return respondFail(c, fiber.StatusUnauthorized, op, "Unauthorized")
		}
		var uid uint
		if claims, ok := tkn.Claims.(jwt.MapClaims); ok {
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
					if n, err1 := strconv.ParseUint(vv, 10, 64); err1 == nil {
						uid = uint(n)
					}
				}
			} else if v, ok := claims["sub"]; ok {
				if s, ok := v.(string); ok {
					if n, err2 := strconv.ParseUint(s, 10, 64); err2 == nil {
						uid = uint(n)
					}
				}
			}
		}
		if uid == 0 {
			return respondFail(c, fiber.StatusUnauthorized, op, "Unauthorized")
		}
		// lookup admin flag from DB
		uRepo := usersRepo.NewRepository(db)
		u, err := uRepo.FindByID(uid)
		if err != nil {
			return respondFail(c, fiber.StatusUnauthorized, op, "Unauthorized")
		}
		isAdmin := false
		if u.IsAdmin != nil {
			isAdmin = *u.IsAdmin
		}
		c.Locals("user_id", uid)
		c.Locals("is_admin", isAdmin)
		return c.Next()
	}
}

// RequireAdmin ensures is_admin is true; otherwise 403 Forbidden with Failed to POST data.
func RequireAdmin(op string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		v := c.Locals("is_admin")
		if v == nil {
			return respondFail(c, fiber.StatusUnauthorized, op, "Unauthorized")
		}
		isAdmin, ok := v.(bool)
		if !ok || !isAdmin {
			return respondFail(c, fiber.StatusForbidden, op, "Forbidden")
		}
		return c.Next()
	}
}

// -------- Handlers --------
// GET /category (public)
func (h *Handler) List(c *fiber.Ctx) error {
	rows, err := h.s.List()
	if err != nil {
		return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
	}
	// map to response
	out := make([]fiber.Map, 0, len(rows))
	for _, it := range rows {
		out = append(out, fiber.Map{"id": it.ID, "nama_category": it.NamaCategory})
	}
	return respondOK(c, "GET", out)
}

// GET /category/:id (public)
func (h *Handler) GetByID(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id64 == 0 {
		return respondFail(c, fiber.StatusNotFound, "GET", "No Data Category")
	}
	cat, err := h.s.GetByID(uint(id64))
	if err != nil {
		// our service/repo returns ErrNotFound as errors.New("not found")
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return respondFail(c, fiber.StatusNotFound, "GET", "No Data Category")
		}
		return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
	}
	return respondOK(c, "GET", fiber.Map{"id": cat.ID, "nama_category": cat.NamaCategory})
}

// POST /category (private admin)
func (h *Handler) Create(c *fiber.Ctx) error {
	var body struct {
		NamaCategory string `json:"nama_category"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondFail(c, fiber.StatusBadRequest, "POST", "Invalid JSON")
	}
	t := strings.TrimSpace(body.NamaCategory)
	if len(t) < 2 {
		return respondFail(c, fiber.StatusBadRequest, "POST", "nama_category minimal 2 karakter")
	}
	isAdmin := false
	if v := c.Locals("is_admin"); v != nil {
		if b, ok := v.(bool); ok {
			isAdmin = b
		}
	}
	id, err := h.s.Create(isAdmin, t)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return respondFail(c, fiber.StatusBadRequest, "POST", "nama_category sudah ada")
		}
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return respondFail(c, fiber.StatusForbidden, "POST", "Forbidden")
		}
		return respondFail(c, fiber.StatusBadRequest, "POST", err.Error())
	}
	return respondOK(c, "POST", int(id))
}

// PUT /category/:id (private admin)
func (h *Handler) Update(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id64 == 0 {
		return respondFail(c, fiber.StatusNotFound, "GET", "No Data Category")
	}
	var body struct {
		NamaCategory string `json:"nama_category"`
	}
	if err := c.BodyParser(&body); err != nil {
		return respondFail(c, fiber.StatusBadRequest, "POST", "Invalid JSON")
	}
	t := strings.TrimSpace(body.NamaCategory)
	if len(t) < 2 {
		return respondFail(c, fiber.StatusBadRequest, "POST", "nama_category minimal 2 karakter")
	}
	isAdmin := false
	if v := c.Locals("is_admin"); v != nil {
		if b, ok := v.(bool); ok {
			isAdmin = b
		}
	}
	if err := h.s.Update(isAdmin, uint(id64), t); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return respondFail(c, fiber.StatusNotFound, "GET", "No Data Category")
		}
		if strings.Contains(strings.ToLower(err.Error()), "duplicate") {
			return respondFail(c, fiber.StatusBadRequest, "POST", "nama_category sudah ada")
		}
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return respondFail(c, fiber.StatusForbidden, "POST", "Forbidden")
		}
		return respondFail(c, fiber.StatusBadRequest, "POST", err.Error())
	}
	return respondOK(c, "GET", "")
}

// DELETE /category/:id (private admin)
func (h *Handler) Delete(c *fiber.Ctx) error {
	idStr := c.Params("id")
	id64, err := strconv.ParseUint(strings.TrimSpace(idStr), 10, 64)
	if err != nil || id64 == 0 {
		// For DELETE invalid ID, treat as record not found -> 400 per requirement
		return respondFail(c, fiber.StatusBadRequest, "GET", "record not found")
	}
	isAdmin := false
	if v := c.Locals("is_admin"); v != nil {
		if b, ok := v.(bool); ok {
			isAdmin = b
		}
	}
	if err := h.s.Delete(isAdmin, uint(id64)); err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return respondFail(c, fiber.StatusBadRequest, "GET", "record not found")
		}
		if strings.Contains(strings.ToLower(err.Error()), "forbidden") {
			return respondFail(c, fiber.StatusForbidden, "POST", "Forbidden")
		}
		return respondFail(c, fiber.StatusBadRequest, "GET", err.Error())
	}
	return respondOK(c, "GET", "")
}
