package auth

import (
	"errors"
	"time"

	"project-evermos/internal/config"
	model "project-evermos/internal/todo/model/auth"
	storemodel "project-evermos/internal/todo/model/toko"
	repo "project-evermos/internal/todo/repository/auth"
	storerepo "project-evermos/internal/todo/repository/toko"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Exported error sentinels for identity comparisons
var (
	ErrInvalidCredentials = errors.New("invalid credentials")
	ErrDuplicateEmail     = errors.New("duplicate_email")
	ErrDuplicatePhone     = errors.New("duplicate_phone")
)

// Service contains business logic for auth domain.
type Service struct {
	repo      *repo.Repository
	storeRepo *storerepo.Repository
	cfg       *config.Config
}

// NewService constructs a new Service.
func NewService(r *repo.Repository, storeR *storerepo.Repository, cfg *config.Config) *Service { return &Service{repo: r, storeRepo: storeR, cfg: cfg} }

// Login authenticates a user by phone and password and returns the user and JWT token.
func (s *Service) Login(phone, password string) (*model.User, string, error) {
	u, err := s.repo.FindByPhone(phone)
	if err != nil {
		return nil, "", err
	}
	if u == nil {
		return nil, "", ErrInvalidCredentials
	}
	if err1 := bcrypt.CompareHashAndPassword([]byte(u.KataSandi), []byte(password)); err1 != nil {
		return nil, "", ErrInvalidCredentials
	}
	// Generate JWT
	exp := time.Now().Add(time.Duration(s.cfg.JWTExpiryDays) * 24 * time.Hour).Unix()
	claims := jwt.MapClaims{
		"email": u.Email,
		"id":    u.ID,
		"exp":   exp,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(s.cfg.JWTSecret))
	if err != nil {
		return nil, "", err
	}
	return u, signed, nil
}

// Register creates a new user after validating uniqueness and hashing the password.
func (s *Service) Register(u *model.User, plainPassword string) error {
	// uniqueness checks
	if existing, err := s.repo.FindByEmail(u.Email); err != nil {
		return err
	} else if existing != nil {
		return ErrDuplicateEmail
	}
	if existing, err := s.repo.FindByPhone(u.NoTelp); err != nil {
		return err
	} else if existing != nil {
		return ErrDuplicatePhone
	}
	// hash password
	h, err := bcrypt.GenerateFromPassword([]byte(plainPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	u.KataSandi = string(h)
	if err := s.repo.CreateUser(u); err != nil {
		return err
	}
	// Auto-create toko for the newly registered user
	newStore := &storemodel.Toko{IDUser: u.ID, NamaToko: u.Nama}
	if err := s.storeRepo.Create(newStore); err != nil {
		return err
	}
	return nil
}
