package model

import (
	"time"
)

// RefreshToken representa la entidad de un token de refresco en la base de datos
type RefreshToken struct {
	ID         string    `json:"id"`
	UserID     string    `json:"user_id"`
	TokenHash  string    `json:"-"` // Hash del token opaco (ej. bcrypt)
	ExpiresAt  time.Time `json:"expires_at"`
	Revoked    bool      `json:"revoked"`
	DeviceInfo string    `json:"device_info"`
	IPAddress  string    `json:"ip_address"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// IsValid verifica si el token no está revocado y no ha expirado
func (r *RefreshToken) IsValid() bool {
	return !r.Revoked && time.Now().Before(r.ExpiresAt)
}
