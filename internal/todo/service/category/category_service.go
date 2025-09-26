package category

import (
	"errors"
	"strings"

	repo "project-evermos/internal/todo/repository/category"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
)

type Service struct{ r *repo.Repository }

func NewService(r *repo.Repository) *Service { return &Service{r: r} }

func (s *Service) List() ([]CategoryDTO, error) {
	rows, err := s.r.List()
	if err != nil { return nil, err }
	out := make([]CategoryDTO, 0, len(rows))
	for _, c := range rows {
		out = append(out, CategoryDTO{ID: c.ID, NamaCategory: c.NamaCategory})
	}
	return out, nil
}

func (s *Service) GetByID(id uint) (*CategoryDTO, error) {
	c, err := s.r.GetByID(id)
	if err != nil { return nil, err }
	return &CategoryDTO{ID: c.ID, NamaCategory: c.NamaCategory}, nil
}

func (s *Service) Create(isAdmin bool, name string) (uint, error) {
	if !isAdmin { return 0, ErrForbidden }
	t := strings.TrimSpace(name)
	if len(t) < 2 { return 0, errors.New("invalid name") }
	id, err := s.r.Create(t)
	if err != nil { return 0, err }
	return id, nil
}

func (s *Service) Update(isAdmin bool, id uint, name string) error {
	if !isAdmin { return ErrForbidden }
	t := strings.TrimSpace(name)
	if len(t) < 2 { return errors.New("invalid name") }
	return s.r.Update(id, t)
}

func (s *Service) Delete(isAdmin bool, id uint) error {
	if !isAdmin { return ErrForbidden }
	return s.r.Delete(id)
}

// DTO used by handlers
type CategoryDTO struct {
	ID           uint   `json:"id"`
	NamaCategory string `json:"nama_category"`
}