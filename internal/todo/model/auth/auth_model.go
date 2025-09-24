package auth

// User represents a simplified user model.
type User struct {
    ID   uint   `gorm:"primaryKey"`
    Name string `gorm:"size:100"`
}