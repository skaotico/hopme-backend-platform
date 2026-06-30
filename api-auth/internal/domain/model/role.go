package model

import "time"

// Role representa un rol en el sistema (por ejemplo: 'SUPER_ADMIN', 'USER')
type Role struct {
	ID          string       `json:"id"`
	Code        string       `json:"code"`
	Name        string       `json:"name"`
	Description string       `json:"description,omitempty"`
	IsSystem    bool         `json:"is_system"`
	Permissions []Permission `json:"permissions,omitempty"`
	CreatedAt   time.Time    `json:"created_at,omitempty"`
}
