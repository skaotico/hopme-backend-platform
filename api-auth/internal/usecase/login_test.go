package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log/slog"
	"os"
	"testing"
	"time"

	"c4-auth/internal/domain/model"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func TestLoginUseCase_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	jwtSecret := "my-secret-key-12345"
	password := "supersecret123"
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)

	activeUser := &model.User{
		ID:           "user-123",
		Username:     "skaotico",
		Email:        "skaotico@homelab.local",
		PasswordHash: string(hashedPassword),
		IsActive:     true,
		Roles: []model.Role{
			{
				Code: "ADMIN",
				Permissions: []model.Permission{
					{Code: "USER_READ"},
				},
			},
		},
	}

	inactiveUser := &model.User{
		ID:           "inactive-123",
		Username:     "inactive",
		Email:        "inactive@homelab.local",
		PasswordHash: string(hashedPassword),
		IsActive:     false,
	}

	tests := []struct {
		name          string
		email         string
		password      string
		ipAddress     string
		deviceInfo    string
		setupMocks    func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository)
		wantErr       bool
		expectedErr   error
		validateToken func(t *testing.T, accessToken, refreshToken string, rtr *mockRefreshTokenRepository)
	}{
		{
			name:       "successful login",
			email:      "skaotico@homelab.local",
			password:   password,
			ipAddress:  "192.168.1.50",
			deviceInfo: "Mozilla/Firefox",
			setupMocks: func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				ur.users["skaotico@homelab.local"] = activeUser
			},
			wantErr: false,
			validateToken: func(t *testing.T, accessToken, refreshToken string, rtr *mockRefreshTokenRepository) {
				if accessToken == "" {
					t.Error("expected access token to be generated, got empty string")
				}
				if refreshToken == "" {
					t.Error("expected refresh token to be generated, got empty string")
				}

				// Check access token claims
				parsedToken, err := jwt.ParseWithClaims(accessToken, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
					return []byte(jwtSecret), nil
				})
				if err != nil {
					t.Fatalf("failed to parse generated access token: %v", err)
				}
				claims, ok := parsedToken.Claims.(*CustomClaims)
				if !ok || !parsedToken.Valid {
					t.Fatal("generated access token is invalid or claims are incorrect")
				}

				if claims.UserID != "user-123" || claims.Username != "skaotico" || claims.Email != "skaotico@homelab.local" {
					t.Errorf("token claims mismatched: %+v", claims)
				}
				if len(claims.Roles) != 1 || claims.Roles[0] != "ADMIN" {
					t.Errorf("expected role ADMIN in token claims, got %v", claims.Roles)
				}
				if len(claims.Permissions) != 1 || claims.Permissions[0] != "USER_READ" {
					t.Errorf("expected permission USER_READ in token claims, got %v", claims.Permissions)
				}

				// Check if the refresh token hash is in the mock db
				hasher := sha256.New()
				hasher.Write([]byte(refreshToken))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				dbRt, err := rtr.FindByHash(context.Background(), rtHash)
				if err != nil || dbRt == nil {
					t.Fatal("expected refresh token to be persisted in database")
				}
				if dbRt.UserID != "user-123" {
					t.Errorf("expected refresh token to belong to user-123, got %s", dbRt.UserID)
				}
				if dbRt.IPAddress != "192.168.1.50" || dbRt.DeviceInfo != "Mozilla/Firefox" {
					t.Errorf("expected metadata to match, got ip:%s device:%s", dbRt.IPAddress, dbRt.DeviceInfo)
				}
			},
		},
		{
			name:      "brute force block by rate limiter",
			email:     "skaotico@homelab.local",
			password:  password,
			ipAddress: "192.168.1.99",
			setupMocks: func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				cr.onIncrAttempts = func(ip string, ttl int) (int, error) {
					return 6, nil // above limit (maxAttempts is 5)
				}
			},
			wantErr:     true,
			expectedErr: model.ErrUnauthorized,
		},
		{
			name:      "user not found",
			email:     "notfound@homelab.local",
			password:  password,
			ipAddress: "192.168.1.50",
			setupMocks: func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				// No user added
			},
			wantErr:     true,
			expectedErr: model.ErrInvalidCredentials,
		},
		{
			name:      "user inactive account",
			email:     "inactive@homelab.local",
			password:  password,
			ipAddress: "192.168.1.50",
			setupMocks: func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				ur.users["inactive@homelab.local"] = inactiveUser
			},
			wantErr:     true,
			expectedErr: model.ErrUnauthorized,
		},
		{
			name:      "wrong password",
			email:     "skaotico@homelab.local",
			password:  "wrongpassword",
			ipAddress: "192.168.1.50",
			setupMocks: func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				ur.users["skaotico@homelab.local"] = activeUser
			},
			wantErr:     true,
			expectedErr: model.ErrInvalidCredentials,
		},
		{
			name:      "database error during refresh token persistence",
			email:     "skaotico@homelab.local",
			password:  password,
			ipAddress: "192.168.1.50",
			setupMocks: func(ur *mockUserRepository, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				ur.users["skaotico@homelab.local"] = activeUser
				rtr.onCreate = func(token *model.RefreshToken) error {
					return errors.New("db write failure")
				}
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := newMockUserRepository()
			cr := newMockCacheRepository()
			rtr := newMockRefreshTokenRepository()

			if tt.setupMocks != nil {
				tt.setupMocks(ur, cr, rtr)
			}

			uc := NewLoginUseCase(
				ur,
				cr,
				rtr,
				jwtSecret,
				15*time.Minute,
				7*24*time.Hour,
				5,
				15,
				logger,
			)

			at, rt, err := uc.Execute(context.Background(), tt.email, tt.password, tt.ipAddress, tt.deviceInfo)

			if (err != nil) != tt.wantErr {
				t.Fatalf("LoginUseCase.Execute() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if tt.expectedErr != nil && !errors.Is(err, tt.expectedErr) {
					t.Errorf("LoginUseCase.Execute() expected error %v, got %v", tt.expectedErr, err)
				}
				return
			}

			if tt.validateToken != nil {
				tt.validateToken(t, at, rt, rtr)
			}
		})
	}
}
