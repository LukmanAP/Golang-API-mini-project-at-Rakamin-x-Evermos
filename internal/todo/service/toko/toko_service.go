package toko

import (
	"errors"
	"strings"
	"time"

	repo "project-evermos/internal/todo/repository/toko"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrNotFound     = errors.New("not_found")
)

// Service contains business logic for toko domain.
type Service struct {
	repo *repo.Repository
}

func NewService(r *repo.Repository) *Service { return &Service{repo: r} }

// GetMyStore returns the store owned by the given userID.
func (s *Service) GetMyStore(userID uint) (map[string]interface{}, error) {
	t, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, nil
	}
	return map[string]interface{}{
		"id":        t.ID,
		"nama_toko": strings.TrimSpace(t.NamaToko),
		"url_foto":  strings.TrimSpace(t.UrlFoto),
	}, nil
}

// UpdateStore updates store by id. Only owner can update.
func (s *Service) UpdateStore(id uint, userID uint, nama string, urlFoto string) error {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return err
	}
	if t == nil {
		return ErrNotFound
	}
	if t.IDUser != userID {
		return ErrForbidden
	}
	if len(strings.TrimSpace(nama)) < 3 {
		return errors.New("nama_toko minimal 3 karakter")
	}
	t.NamaToko = strings.TrimSpace(nama)
	t.UrlFoto = strings.TrimSpace(urlFoto)
	t.UpdatedAt = time.Now()
	return s.repo.Update(t)
}

// GetByID returns store by id. If public access, omit user_id.
func (s *Service) GetByID(id uint, public bool, requesterUserID uint) (map[string]interface{}, error) {
	t, err := s.repo.FindByID(id)
	if err != nil {
		return nil, err
	}
	if t == nil {
		return nil, ErrNotFound
	}
	if !public && t.IDUser != requesterUserID {
		return nil, ErrForbidden
	}
	resp := map[string]interface{}{
		"id":        t.ID,
		"nama_toko": strings.TrimSpace(t.NamaToko),
		"url_foto":  strings.TrimSpace(t.UrlFoto),
	}
	return resp, nil
}

// List returns paginated stores.
func (s *Service) List(limit, page int, name string) (map[string]interface{}, error) {
	if limit <= 0 {
		limit = 10
	}
	if page <= 0 {
		page = 1
	}
	items, total, err := s.repo.List(limit, page, name)
	if err != nil {
		return nil, err
	}
	respItems := make([]map[string]interface{}, 0, len(items))
	for _, t := range items {
		respItems = append(respItems, map[string]interface{}{
			"id":        t.ID,
			"nama_toko": strings.TrimSpace(t.NamaToko),
			"url_foto":  strings.TrimSpace(t.UrlFoto),
		})
	}
	return map[string]interface{}{
		"items":      respItems,
		"total":      total,
		"page":       page,
		"limit":      limit,
		"total_page": (total + int64(limit) - 1) / int64(limit),
	}, nil
}