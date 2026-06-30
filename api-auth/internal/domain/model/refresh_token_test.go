package model

import (
	"testing"
	"time"
)

func TestRefreshToken_IsValid(t *testing.T) {
	tests := []struct {
		name         string
		refreshToken RefreshToken
		want         bool
	}{
		{
			name: "token is valid (not revoked, not expired)",
			refreshToken: RefreshToken{
				Revoked:   false,
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			want: true,
		},
		{
			name: "token is invalid - revoked is true",
			refreshToken: RefreshToken{
				Revoked:   true,
				ExpiresAt: time.Now().Add(1 * time.Hour),
			},
			want: false,
		},
		{
			name: "token is invalid - expired in the past",
			refreshToken: RefreshToken{
				Revoked:   false,
				ExpiresAt: time.Now().Add(-1 * time.Hour),
			},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := tt.refreshToken.IsValid()
			if got != tt.want {
				t.Errorf("RefreshToken.IsValid() = %v, want = %v", got, tt.want)
			}
		})
	}
}
