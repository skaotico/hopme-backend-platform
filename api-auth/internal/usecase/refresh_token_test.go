package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"testing"
	"time"

	"c4-auth/internal/domain/model"
)

func TestRefreshTokenUseCase_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	jwtSecret := "my-secret-key-12345"

	activeUser := &model.User{
		ID:       "user-123",
		Username: "skaotico",
		Email:    "skaotico@homelab.local",
		IsActive: true,
	}

	inactiveUser := &model.User{
		ID:       "user-inactive",
		Username: "inactive",
		Email:    "inactive@homelab.local",
		IsActive: false,
	}

	tests := []struct {
		name         string
		refreshToken string
		ipAddress    string
		deviceInfo   string
		setupMocks   func(ur *mockUserRepository, rtr *mockRefreshTokenRepository)
		wantErr      bool
		expectedErr  string
		validate     func(t *testing.T, at, rt string, ur *mockUserRepository, rtr *mockRefreshTokenRepository)
	}{
		{
			name:         "successful token rotation",
			refreshToken: "original-plain-rt",
			ipAddress:    "192.168.1.100",
			deviceInfo:   "Mozilla/Chrome",
			setupMocks: func(ur *mockUserRepository, rtr *mockRefreshTokenRepository) {
				ur.usersByID["user-123"] = activeUser

				hasher := sha256.New()
				hasher.Write([]byte("original-plain-rt"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				originalRt := &model.RefreshToken{
					ID:        "rt-uuid-1",
					UserID:    "user-123",
					TokenHash: rtHash,
					ExpiresAt: time.Now().Add(1 * time.Hour),
					Revoked:   false,
				}
				rtr.tokens[rtHash] = originalRt
				rtr.tokensByID["rt-uuid-1"] = originalRt
			},
			wantErr: false,
			validate: func(t *testing.T, at, rt string, ur *mockUserRepository, rtr *mockRefreshTokenRepository) {
				if at == "" || rt == "" {
					t.Fatal("expected access and refresh tokens to be returned")
				}

				// Verify old token is revoked
				hasher1 := sha256.New()
				hasher1.Write([]byte("original-plain-rt"))
				oldHash := hex.EncodeToString(hasher1.Sum(nil))
				if !rtr.tokens[oldHash].Revoked {
					t.Error("expected old refresh token to be revoked upon rotation")
				}

				// Verify new token is created
				hasher2 := sha256.New()
				hasher2.Write([]byte(rt))
				newHash := hex.EncodeToString(hasher2.Sum(nil))

				newRt := rtr.tokens[newHash]
				if newRt == nil {
					t.Fatal("expected new refresh token to be saved in DB")
				}
				if newRt.UserID != "user-123" {
					t.Errorf("expected new token user ID to be 'user-123', got '%s'", newRt.UserID)
				}
				if newRt.IPAddress != "192.168.1.100" || newRt.DeviceInfo != "Mozilla/Chrome" {
					t.Errorf("expected new token metadata to be updated, got ip:%s device:%s", newRt.IPAddress, newRt.DeviceInfo)
				}
			},
		},
		{
			name:         "token already revoked - reuse security alert",
			refreshToken: "revoked-plain-rt",
			setupMocks: func(ur *mockUserRepository, rtr *mockRefreshTokenRepository) {
				hasher := sha256.New()
				hasher.Write([]byte("revoked-plain-rt"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				revokedRt := &model.RefreshToken{
					ID:        "rt-uuid-revoked",
					UserID:    "user-123",
					TokenHash: rtHash,
					ExpiresAt: time.Now().Add(1 * time.Hour),
					Revoked:   true,
				}
				rtr.tokens[rtHash] = revokedRt
				rtr.tokensByID["rt-uuid-revoked"] = revokedRt

				// Add another token belonging to same user to test family revocation
				otherHasher := sha256.New()
				otherHasher.Write([]byte("other-rt"))
				otherHash := hex.EncodeToString(otherHasher.Sum(nil))
				otherRt := &model.RefreshToken{
					ID:        "rt-uuid-other",
					UserID:    "user-123",
					TokenHash: otherHash,
					ExpiresAt: time.Now().Add(1 * time.Hour),
					Revoked:   false,
				}
				rtr.tokens[otherHash] = otherRt
				rtr.tokensByID["rt-uuid-other"] = otherRt
			},
			wantErr:     true,
			expectedErr: "token reused - security alert",
			validate: func(t *testing.T, at, rt string, ur *mockUserRepository, rtr *mockRefreshTokenRepository) {
				// Verify whole family got revoked
				for id, token := range rtr.tokensByID {
					if token.UserID == "user-123" && !token.Revoked {
						t.Errorf("expected token %s belonging to user-123 to be revoked as part of security alert", id)
					}
				}
			},
		},
		{
			name:         "token expired",
			refreshToken: "expired-plain-rt",
			setupMocks: func(ur *mockUserRepository, rtr *mockRefreshTokenRepository) {
				hasher := sha256.New()
				hasher.Write([]byte("expired-plain-rt"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				expiredRt := &model.RefreshToken{
					ID:        "rt-uuid-expired",
					UserID:    "user-123",
					TokenHash: rtHash,
					ExpiresAt: time.Now().Add(-5 * time.Minute), // expired in past
					Revoked:   false,
				}
				rtr.tokens[rtHash] = expiredRt
				rtr.tokensByID["rt-uuid-expired"] = expiredRt
			},
			wantErr:     true,
			expectedErr: "refresh token expired",
		},
		{
			name:         "token not found",
			refreshToken: "non-existent-rt",
			wantErr:      true,
			expectedErr:  "invalid refresh token",
		},
		{
			name:         "user inactive",
			refreshToken: "active-rt-inactive-user",
			setupMocks: func(ur *mockUserRepository, rtr *mockRefreshTokenRepository) {
				ur.usersByID["user-inactive"] = inactiveUser

				hasher := sha256.New()
				hasher.Write([]byte("active-rt-inactive-user"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				rt := &model.RefreshToken{
					ID:        "rt-uuid-inactive",
					UserID:    "user-inactive",
					TokenHash: rtHash,
					ExpiresAt: time.Now().Add(1 * time.Hour),
					Revoked:   false,
				}
				rtr.tokens[rtHash] = rt
				rtr.tokensByID["rt-uuid-inactive"] = rt
			},
			wantErr:     true,
			expectedErr: "user not found or inactive",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ur := newMockUserRepository()
			rtr := newMockRefreshTokenRepository()

			if tt.setupMocks != nil {
				tt.setupMocks(ur, rtr)
			}

			uc := NewRefreshTokenUseCase(
				rtr,
				ur,
				jwtSecret,
				15*time.Minute,
				7*24*time.Hour,
				logger,
			)

			at, rt, err := uc.Execute(context.Background(), tt.refreshToken, tt.ipAddress, tt.deviceInfo)

			if (err != nil) != tt.wantErr {
				t.Fatalf("RefreshTokenUseCase.Execute() error = %v, wantErr = %v", err, tt.wantErr)
			}

			if tt.wantErr {
				if tt.expectedErr != "" && (err == nil || err.Error() != tt.expectedErr) {
					t.Errorf("RefreshTokenUseCase.Execute() expected error '%s', got %v", tt.expectedErr, err)
				}
				return
			}

			if tt.validate != nil {
				tt.validate(t, at, rt, ur, rtr)
			}
		})
	}
}
