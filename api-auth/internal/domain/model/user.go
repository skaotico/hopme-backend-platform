package model

import (
	"errors"
	"regexp"
	"strings"
	"time"
)

var (
	ErrInvalidEmail    = errors.New("el formato del correo electrónico es inválido")
	ErrInvalidUsername = errors.New("el nombre de usuario debe tener entre 3 y 50 caracteres y no contener caracteres especiales")
	ErrInvalidPassword = errors.New("la contraseña debe tener al menos 6 caracteres")
)

var emailRegex = regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,15}$`)

// User representa la entidad Usuario del dominio (esquema auth.users)
type User struct {
	ID           string     `json:"id"`
	Username     string     `json:"username"`
	Email        string     `json:"email"`
	PasswordHash string     `json:"-"`
	DisplayName  string     `json:"display_name,omitempty"`
	AvatarURL    string     `json:"avatar_url,omitempty"`
	IsActive     bool       `json:"is_active"`
	IsVerified   bool       `json:"is_verified"`
	LastLoginAt  *time.Time `json:"last_login_at,omitempty"`
	Roles        []Role     `json:"roles,omitempty"`
	CreatedAt    time.Time  `json:"created_at"`
	UpdatedAt    time.Time  `json:"updated_at"`
}

// Validate realiza validaciones semánticas puras en la entidad de dominio
func (u *User) Validate() error {
	u.Username = strings.TrimSpace(u.Username)
	u.Email = strings.TrimSpace(strings.ToLower(u.Email))
	u.DisplayName = strings.TrimSpace(u.DisplayName)

	if len(u.Username) < 3 || len(u.Username) > 50 {
		return ErrInvalidUsername
	}

	if !emailRegex.MatchString(u.Email) {
		return ErrInvalidEmail
	}

	return nil
}

// HasRole verifica si el usuario posee un rol específico basado en su código
func (u *User) HasRole(roleCode string) bool {
	for _, role := range u.Roles {
		if strings.EqualFold(role.Code, roleCode) {
			return true
		}
	}
	return false
}

// HasPermission verifica si el usuario posee un permiso granular específico
func (u *User) HasPermission(permissionCode string) bool {
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			if strings.EqualFold(perm.Code, permissionCode) {
				return true
			}
		}
	}
	return false
}

// GetAllPermissions recopila todos los códigos de permisos del usuario sin duplicados
func (u *User) GetAllPermissions() []string {
	permMap := make(map[string]bool)
	for _, role := range u.Roles {
		for _, perm := range role.Permissions {
			permMap[perm.Code] = true
		}
	}

	perms := make([]string, 0, len(permMap))
	for perm := range permMap {
		perms = append(perms, perm)
	}
	return perms
}
