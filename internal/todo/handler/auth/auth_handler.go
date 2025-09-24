package auth

// Handler handles auth-related HTTP requests.
type Handler struct{}

// NewHandler constructs a new Handler.
func NewHandler() *Handler { return &Handler{} }