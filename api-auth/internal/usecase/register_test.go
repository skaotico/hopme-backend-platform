package usecase

import (
	"context"
	"errors"
	"log/slog"
	"os"
	"testing"

	"c4-auth/internal/domain/model"
)

func TestRegisterUseCase_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name         string
		username     string
		email        string
		password     string
		setupMocks   func(ur *mockUserRepository)
		wantErr      bool
		expectedErr  error
		validateUser func(t *testing.T, u *model.User)
	}{
		{
			name:     "successful registration",
			username: "skaotico",
			email:    "skaotico@homelab.local",
			password: "secretpassword123",
			setupMocks: func(ur *mockUserRepository) {
				ur.roles["USER"] = &model.Role{
					Code: "USER",
					Name: "Usuario estándar",
				}
			},
			wantErr: false,
			validateUser: func(t *testing.T, u *model.User) {
				if u == nil {
					t.Fatal("expected user to be returned, got nil")
				}
				if u.Username != "skaotico" {
					t.Errorf("expected username to be 'skaotico', got '%s'", u.Username)
				}
				if u.Email != "skaotico@homelab.local" {
					t.Errorf("expected email to be 'skaotico@homelab.local', got '%s'", u.Email)
				}
				if u.PasswordHash == "" || u.PasswordHash == "secretpassword123" {
					t.Error("expected password hash to be generated and not plain")
				}
				if len(u.Roles) != 1 || u.Roles[0].Code != "USER" {
					t.Errorf("expected user to have the default 'USER' role, got %v", u.Roles)
				}
			},
		},
		{
			name:     "invalid username (too short)",
			username: "sk",
			email:    "skaotico@homelab.local",
			password: "secretpassword123",
			wantErr:  true,
			expectedErr: model.ErrInvalidUsername,
		},
		{
			name:     "invalid email format",
			username: "skaotico",
			email:    "skaotico.homelab.local",
			password: "secretpassword123",
			wantErr:  true,
			expectedErr: model.ErrInvalidEmail,
		},
		{
			name:     "password too short",
			username: "skaotico",
			email:    "skaotico@homelab.local",
			password: "123",
			wantErr:  true,
			expectedErr: model.ErrInvalidPassword,
		},
		{
			name:     "email already exists",
			username: "skaotico",
			email:    "skaotico@homelab.local",
			password: "secretpassword123",
			setupMocks: func(ur *mockUserRepository) {
				ur.users["skaotico@homelab.local"] = &model.User{
					ID:       "existing-id",
					Username: "someother",
					Email:    "skaotico@homelab.local",
				}
			},
			wantErr:     true,
			expectedErr: model.ErrUserAlreadyExists,
		},
		{
			name:     "username already exists",
			username: "skaotico",
			email:    "skaotico@homelab.local",
			password: "secretpassword123",
			setupMocks: func(ur *mockUserRepository) {
				ur.usersByUsername["skaotico"] = &model.User{
					ID:       "existing-id",
					Username: "skaotico",
					Email:    "other@homelab.local",
				}
			},
			wantErr:     true,
			expectedErr: model.ErrUserAlreadyExists,
		},
		{
			name:     "database error on create",
			username: "skaotico",
			email:    "skaotico@homelab.local",
			password: "secretpassword123",
			setupMocks: func(ur *mockUserRepository) {
				ur.onCreate = func(user *model.User) error {
					return errors.New("db connection failure")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := newMockUserRepository()
			if tt.setupMocks != nil {
				tt.setupMocks(ur)
			}

			uc := NewRegisterUseCase(ur, logger)
			user, err := uc.Execute(context.Background(), tt.username, tt.email, tt.password)

			if (err != nil) != tt.wantErr {
				t.Fatalf("RegisterUseCase.Execute() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("RegisterUseCase.Execute() expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if tt.validateUser != nil {
				tt.validateUser(t, user)
			}
		})
	}
}
