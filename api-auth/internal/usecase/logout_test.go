package usecase

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"log/slog"
	"os"
	"testing"

	"c4-auth/internal/domain/model"
)

func TestLogoutUseCase_Execute(t *testing.T) {
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))

	tests := []struct {
		name         string
		jti          string
		timeToLive   int
		refreshToken string
		setupMocks   func(cr *mockCacheRepository, rtr *mockRefreshTokenRepository)
		validate     func(t *testing.T, cr *mockCacheRepository, rtr *mockRefreshTokenRepository)
	}{
		{
			name:         "logout with both jwt and refresh token successfully",
			jti:          "jwt-jti-111",
			timeToLive:   300,
			refreshToken: "plain-rt-222",
			setupMocks: func(cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				hasher := sha256.New()
				hasher.Write([]byte("plain-rt-222"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				rtr.tokens[rtHash] = &model.RefreshToken{
					ID:        "token-uuid",
					TokenHash: rtHash,
					UserID:    "user-123",
					Revoked:   false,
				}
				rtr.tokensByID["token-uuid"] = rtr.tokens[rtHash]
			},
			validate: func(t *testing.T, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				// Check JWT blacklisted
				blacklisted, _ := cr.IsJWTBlacklisted(context.Background(), "jwt-jti-111")
				if !blacklisted {
					t.Error("expected JWT to be blacklisted in cache")
				}

				// Check Refresh Token revoked
				hasher := sha256.New()
				hasher.Write([]byte("plain-rt-222"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				rt := rtr.tokens[rtHash]
				if rt == nil || !rt.Revoked {
					t.Error("expected refresh token to be revoked in database")
				}
			},
		},
		{
			name:         "logout with only jwt successfully",
			jti:          "jwt-jti-only",
			timeToLive:   120,
			refreshToken: "",
			validate: func(t *testing.T, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				blacklisted, _ := cr.IsJWTBlacklisted(context.Background(), "jwt-jti-only")
				if !blacklisted {
					t.Error("expected JWT to be blacklisted in cache")
				}
			},
		},
		{
			name:         "logout with only refresh token successfully",
			jti:          "",
			timeToLive:   0,
			refreshToken: "plain-rt-only",
			setupMocks: func(cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				hasher := sha256.New()
				hasher.Write([]byte("plain-rt-only"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				rtr.tokens[rtHash] = &model.RefreshToken{
					ID:        "token-uuid-only",
					TokenHash: rtHash,
					UserID:    "user-123",
					Revoked:   false,
				}
				rtr.tokensByID["token-uuid-only"] = rtr.tokens[rtHash]
			},
			validate: func(t *testing.T, cr *mockCacheRepository, rtr *mockRefreshTokenRepository) {
				// Blacklist should be untouched
				if len(cr.blacklist) > 0 {
					t.Error("expected blacklist to remain empty")
				}

				// Refresh token should be revoked
				hasher := sha256.New()
				hasher.Write([]byte("plain-rt-only"))
				rtHash := hex.EncodeToString(hasher.Sum(nil))

				rt := rtr.tokens[rtHash]
				if rt == nil || !rt.Revoked {
					t.Error("expected refresh token to be revoked in database")
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cr := newMockCacheRepository()
			rtr := newMockRefreshTokenRepository()

			if tt.setupMocks != nil {
				tt.setupMocks(cr, rtr)
			}

			uc := NewLogoutUseCase(cr, rtr, logger)
			err := uc.Execute(context.Background(), tt.jti, tt.timeToLive, tt.refreshToken)

			if err != nil {
				t.Fatalf("LogoutUseCase.Execute() returned unexpected error: %v", err)
			}

			if tt.validate != nil {
				tt.validate(t, cr, rtr)
			}
		})
	}
}
