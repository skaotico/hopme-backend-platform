package model

import (
	"testing"
)

func TestUser_Validate(t *testing.T) {
	tests := []struct {
		name    string
		user    User
		wantErr error
	}{
		{
			name: "valid user",
			user: User{
				Username: "skaotico",
				Email:    "skaotico@homelab.local",
			},
			wantErr: nil,
		},
		{
			name: "username too short",
			user: User{
				Username: "sk",
				Email:    "skaotico@homelab.local",
			},
			wantErr: ErrInvalidUsername,
		},
		{
			name: "username too long",
			user: User{
				Username: "skaoticoskaoticoskaoticoskaoticoskaoticoskaoticoskaotico",
				Email:    "skaotico@homelab.local",
			},
			wantErr: ErrInvalidUsername,
		},
		{
			name: "invalid email format - no at",
			user: User{
				Username: "skaotico",
				Email:    "skaoticohomelab.local",
			},
			wantErr: ErrInvalidEmail,
		},
		{
			name: "invalid email format - no domain",
			user: User{
				Username: "skaotico",
				Email:    "skaotico@homelab",
			},
			wantErr: ErrInvalidEmail,
		},
		{
			name: "valid email - uppercase should be standardized",
			user: User{
				Username: "skaotico",
				Email:    "SKAOTICO@HOMELAB.LOCAL",
			},
			wantErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.user.Validate()
			if err != tt.wantErr {
				t.Errorf("User.Validate() error = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_HasRole(t *testing.T) {
	roles := []Role{
		{Code: "ADMIN", Name: "Administrador"},
		{Code: "USER", Name: "Usuario estándar"},
	}
	user := User{
		Username: "skaotico",
		Roles:    roles,
	}

	t.Run("matching role code case-insensitive", func(t *testing.T) {
		if !user.HasRole("admin") {
			t.Error("expected HasRole('admin') to be true")
		}
		if !user.HasRole("USER") {
			t.Error("expected HasRole('USER') to be true")
		}
	})

	t.Run("non-matching role code", func(t *testing.T) {
		if user.HasRole("SUPER_ADMIN") {
			t.Error("expected HasRole('SUPER_ADMIN') to be false")
		}
	})

	t.Run("empty roles list", func(t *testing.T) {
		emptyUser := User{Username: "guest"}
		if emptyUser.HasRole("USER") {
			t.Error("expected empty user to not have role USER")
		}
	})
}

func TestUser_HasPermission(t *testing.T) {
	roles := []Role{
		{
			Code: "ADMIN",
			Permissions: []Permission{
				{Code: "USER_READ"},
				{Code: "USER_CREATE"},
			},
		},
		{
			Code: "USER",
			Permissions: []Permission{
				{Code: "MODULE_READ"},
			},
		},
	}
	user := User{
		Username: "skaotico",
		Roles:    roles,
	}

	t.Run("has permission case-insensitive", func(t *testing.T) {
		if !user.HasPermission("user_read") {
			t.Error("expected HasPermission('user_read') to be true")
		}
		if !user.HasPermission("MODULE_READ") {
			t.Error("expected HasPermission('MODULE_READ') to be true")
		}
	})

	t.Run("does not have permission", func(t *testing.T) {
		if user.HasPermission("USER_DELETE") {
			t.Error("expected HasPermission('USER_DELETE') to be false")
		}
	})
}

func TestUser_GetAllPermissions(t *testing.T) {
	roles := []Role{
		{
			Code: "ADMIN",
			Permissions: []Permission{
				{Code: "USER_READ"},
				{Code: "USER_CREATE"},
			},
		},
		{
			Code: "USER",
			Permissions: []Permission{
				{Code: "USER_READ"}, // Duplicate permission
				{Code: "MODULE_READ"},
			},
		},
	}
	user := User{
		Username: "skaotico",
		Roles:    roles,
	}

	perms := user.GetAllPermissions()

	// Verify size and duplicates removal
	if len(perms) != 3 {
		t.Errorf("expected 3 unique permissions, got %d: %v", len(perms), perms)
	}

	permMap := make(map[string]bool)
	for _, p := range perms {
		permMap[p] = true
	}

	expected := []string{"USER_READ", "USER_CREATE", "MODULE_READ"}
	for _, exp := range expected {
		if !permMap[exp] {
			t.Errorf("expected permission list to contain %s", exp)
		}
	}
}
