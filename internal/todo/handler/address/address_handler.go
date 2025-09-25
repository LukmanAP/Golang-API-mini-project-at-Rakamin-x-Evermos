package address

import (
	"context"
	"errors"
	"regexp"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	addrmodel "project-evermos/internal/todo/model/address"
	svc "project-evermos/internal/todo/service/address"
)

type Handler struct{ s *svc.Service }

func NewHandler(s *svc.Service) *Handler { return &Handler{s: s} }

var numRe = regexp.MustCompile(`^\d+$`)

func respondOK(c *fiber.Ctx, data interface{}) error {
	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"status":  true,
		"message": "Succeed to get data",
		"errors":  nil,
		"data":    data,
	})
}

func respondFail(c *fiber.Ctx, code int, msg string) error {
	return c.Status(code).JSON(fiber.Map{
		"status":  false,
		"message": "Failed to get data",
		"errors":  []string{msg},
		"data":    nil,
	})
}

// GET /provcity/listprovincies?search=&limit=&page=
func (h *Handler) ListProvinces(c *fiber.Ctx) error {
	search := c.Query("search", "")
	limitStr := c.Query("limit", "")
	pageStr := c.Query("page", "")
	limit := 0
	page := 1
	if limitStr != "" {
		if v, err := strconv.Atoi(limitStr); err == nil {
			if v < 0 {
				v = 0
			}
			if v > 100 {
				v = 100
			}
			limit = v
		}
	}
	if pageStr != "" {
		if v, err := strconv.Atoi(pageStr); err == nil && v > 0 {
			page = v
		}
	}

	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	items, err := h.s.ListProvinces(ctx, search, limit, page)
	if err != nil {
		switch {
		case errors.Is(err, svc.ErrTimeout):
			return respondFail(c, fiber.StatusGatewayTimeout, "Upstream EMSIFA error")
		case errors.Is(err, svc.ErrUpstream):
			return respondFail(c, fiber.StatusBadGateway, "Upstream EMSIFA error")
		default:
			return respondFail(c, fiber.StatusBadGateway, "Upstream EMSIFA error")
		}
	}
	return respondOK(c, items)
}

// GET /provcity/listcities/:prov_id
func (h *Handler) ListCities(c *fiber.Ctx) error {
	provID := c.Params("prov_id")
	if !numRe.MatchString(provID) {
		return respondFail(c, fiber.StatusBadRequest, "Invalid ID")
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	items, err := h.s.ListCities(ctx, provID)
	if err != nil {
		switch {
		case errors.Is(err, svc.ErrNotFound):
			return respondFail(c, fiber.StatusNotFound, "Data not found")
		case errors.Is(err, svc.ErrTimeout):
			return respondFail(c, fiber.StatusGatewayTimeout, "Upstream EMSIFA error")
		default:
			return respondFail(c, fiber.StatusBadGateway, "Upstream EMSIFA error")
		}
	}
	return respondOK(c, items)
}

// GET /provcity/detailprovince/:prov_id
func (h *Handler) DetailProvince(c *fiber.Ctx) error {
	id := c.Params("prov_id")
	if !numRe.MatchString(id) {
		return respondFail(c, fiber.StatusBadRequest, "Invalid ID")
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	p, err := h.s.DetailProvince(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, svc.ErrNotFound):
			return respondFail(c, fiber.StatusNotFound, "Data not found")
		case errors.Is(err, svc.ErrTimeout):
			return respondFail(c, fiber.StatusGatewayTimeout, "Upstream EMSIFA error")
		default:
			return respondFail(c, fiber.StatusBadGateway, "Upstream EMSIFA error")
		}
	}
	return respondOK(c, addrmodel.Province{ID: p.ID, Name: p.Name})
}

// GET /provcity/detailcity/:city_id
func (h *Handler) DetailCity(c *fiber.Ctx) error {
	id := c.Params("city_id")
	if !numRe.MatchString(id) {
		return respondFail(c, fiber.StatusBadRequest, "Invalid ID")
	}
	ctx, cancel := context.WithTimeout(c.Context(), 5*time.Second)
	defer cancel()
	r, err := h.s.DetailCity(ctx, id)
	if err != nil {
		switch {
		case errors.Is(err, svc.ErrNotFound):
			return respondFail(c, fiber.StatusNotFound, "Data not found")
		case errors.Is(err, svc.ErrTimeout):
			return respondFail(c, fiber.StatusGatewayTimeout, "Upstream EMSIFA error")
		default:
			return respondFail(c, fiber.StatusBadGateway, "Upstream EMSIFA error")
		}
	}
	return respondOK(c, addrmodel.Regency{ID: r.ID, ProvinceID: r.ProvinceID, Name: r.Name})
}