package transaction

import (
    "strconv"

    svc "project-evermos/internal/todo/service/transaction"

    "github.com/gofiber/fiber/v2"
)

type Handler struct { svc *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{svc: s} }

// Response helpers to keep consistent format
func respondOK(c *fiber.Ctx, verb string, data interface{}) error {
    return c.Status(fiber.StatusOK).JSON(fiber.Map{
        "status":  true,
        "message": "Succeed to " + verb + " data",
        "errors":  nil,
        "data":    data,
    })
}

func respondFail(c *fiber.Ctx, code int, verb string, errs []string) error {
    return c.Status(code).JSON(fiber.Map{
        "status":  false,
        "message": "Failed to " + verb + " data",
        "errors":  errs,
        "data":    nil,
    })
}

func jwtUserID(c *fiber.Ctx) (uint, bool) {
    v := c.Locals("user_id")
    if v == nil { return 0, false }
    switch t := v.(type) {
    case int: return uint(t), true
    case int32: return uint(t), true
    case int64: return uint(t), true
    case uint: return t, true
    case uint32: return uint(t), true
    case uint64: return uint(t), true
    case float64: return uint(t), true
    default:
        return 0, false
    }
}

// GET /trx
func (h *Handler) List(c *fiber.Ctx) error {
    uid, ok := jwtUserID(c)
    if !ok { return respondFail(c, fiber.StatusUnauthorized, "GET", []string{"Unauthorized"}) }

    limit, _ := strconv.Atoi(c.Query("limit", "10"))
    page, _ := strconv.Atoi(c.Query("page", "1"))

    resp, err := h.svc.List(uid, limit, page)
    if err != nil { return respondFail(c, fiber.StatusBadRequest, "GET", []string{err.Error()}) }

    return respondOK(c, "GET", resp)
}

// GET /trx/:id
func (h *Handler) GetByID(c *fiber.Ctx) error {
    uid, ok := jwtUserID(c)
    if !ok { return respondFail(c, fiber.StatusUnauthorized, "GET", []string{"Unauthorized"}) }

    id64, _ := strconv.ParseUint(c.Params("id"), 10, 64)
    if id64 == 0 { return respondFail(c, fiber.StatusBadRequest, "GET", []string{"invalid id"}) }

    item, err := h.svc.GetByID(uint(id64), uid)
    if err != nil {
        switch err.Error() {
        case "forbidden":
            return respondFail(c, fiber.StatusForbidden, "GET", []string{"Forbidden"})
        case "not found":
            return respondFail(c, fiber.StatusNotFound, "GET", []string{"No Data Trx"})
        default:
            return respondFail(c, fiber.StatusBadRequest, "GET", []string{err.Error()})
        }
    }

    return respondOK(c, "GET", item)
}

// POST /trx
func (h *Handler) Create(c *fiber.Ctx) error {
    uid, ok := jwtUserID(c)
    if !ok { return respondFail(c, fiber.StatusUnauthorized, "POST", []string{"Unauthorized"}) }

    var req svc.CreateRequest
    if err := c.BodyParser(&req); err != nil {
        return respondFail(c, fiber.StatusBadRequest, "POST", []string{"invalid payload"})
    }

    id, err := h.svc.Create(uid, req)
    if err != nil {
        switch err.Error() {
        case "alamat not owned by user":
            return respondFail(c, fiber.StatusForbidden, "POST", []string{"Alamat bukan milik user"})
        case "alamat not found":
            return respondFail(c, fiber.StatusBadRequest, "POST", []string{"Alamat tidak ditemukan"})
        case "product not found":
            return respondFail(c, fiber.StatusBadRequest, "POST", []string{"Product tidak valid"})
        default:
            return respondFail(c, fiber.StatusBadRequest, "POST", []string{err.Error()})
        }
    }

    return respondOK(c, "POST", id)
}