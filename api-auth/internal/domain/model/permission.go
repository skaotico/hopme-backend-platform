package model

import "time"

// Permission representa un permiso granular dentro del esquema de IAM
type Permission struct {
	ID          string    `json:"id"`
	Code        string    `json:"code"`
	Description string    `json:"description,omitempty"`
	CreatedAt   time.Time `json:"created_at,omitempty"`
}
