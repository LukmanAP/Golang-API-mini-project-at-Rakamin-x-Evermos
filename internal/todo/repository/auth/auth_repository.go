package auth

// Repository handles data access for auth domain.
type Repository struct{}

// NewRepository constructs a new Repository.
func NewRepository() *Repository { return &Repository{} }