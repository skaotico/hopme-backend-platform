package usecase

import (
	"context"
	"c4-auth/internal/domain/model"
)

type mockUserRepository struct {
	users           map[string]*model.User // by email
	usersByUsername map[string]*model.User
	usersByID      map[string]*model.User
	roles          map[string]*model.Role
	onCreate       func(user *model.User) error
	onAssignRole   func(userID, roleCode string) error
}

func newMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users:           make(map[string]*model.User),
		usersByUsername: make(map[string]*model.User),
		usersByID:      make(map[string]*model.User),
		roles:           make(map[string]*model.Role),
	}
}

func (m *mockUserRepository) Create(ctx context.Context, user *model.User) error {
	if m.onCreate != nil {
		return m.onCreate(user)
	}
	user.ID = "generated-user-id"
	m.users[user.Email] = user
	m.usersByUsername[user.Username] = user
	m.usersByID[user.ID] = user
	return nil
}

func (m *mockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	u, ok := m.users[email]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	u, ok := m.usersByUsername[username]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	u, ok := m.usersByID[id]
	if !ok {
		return nil, nil
	}
	return u, nil
}

func (m *mockUserRepository) AssignRole(ctx context.Context, userID, roleCode string) error {
	if m.onAssignRole != nil {
		return m.onAssignRole(userID, roleCode)
	}
	return nil
}

func (m *mockUserRepository) FindRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	r, ok := m.roles[code]
	if !ok {
		return nil, nil
	}
	return r, nil
}

type mockRefreshTokenRepository struct {
	tokens         map[string]*model.RefreshToken // by hash
	tokensByID     map[string]*model.RefreshToken
	onCreate       func(token *model.RefreshToken) error
	onRevoke       func(id string) error
	onRevokeFamily func(userID string) error
}

func newMockRefreshTokenRepository() *mockRefreshTokenRepository {
	return &mockRefreshTokenRepository{
		tokens:     make(map[string]*model.RefreshToken),
		tokensByID: make(map[string]*model.RefreshToken),
	}
}

func (m *mockRefreshTokenRepository) Create(ctx context.Context, token *model.RefreshToken) error {
	if m.onCreate != nil {
		return m.onCreate(token)
	}
	token.ID = "generated-rt-id"
	m.tokens[token.TokenHash] = token
	m.tokensByID[token.ID] = token
	return nil
}

func (m *mockRefreshTokenRepository) FindByHash(ctx context.Context, hash string) (*model.RefreshToken, error) {
	t, ok := m.tokens[hash]
	if !ok {
		return nil, nil
	}
	return t, nil
}

func (m *mockRefreshTokenRepository) Revoke(ctx context.Context, id string) error {
	if m.onRevoke != nil {
		return m.onRevoke(id)
	}
	if t, ok := m.tokensByID[id]; ok {
		t.Revoked = true
	}
	return nil
}

func (m *mockRefreshTokenRepository) RevokeFamily(ctx context.Context, userID string) error {
	if m.onRevokeFamily != nil {
		return m.onRevokeFamily(userID)
	}
	for _, t := range m.tokens {
		if t.UserID == userID {
			t.Revoked = true
		}
	}
	return nil
}

type mockCacheRepository struct {
	blacklist      map[string]bool
	attempts       map[string]int
	profiles       map[string]string
	onIncrAttempts func(ip string, ttl int) (int, error)
}

func newMockCacheRepository() *mockCacheRepository {
	return &mockCacheRepository{
		blacklist: make(map[string]bool),
		attempts:  make(map[string]int),
		profiles:  make(map[string]string),
	}
}

func (m *mockCacheRepository) BlacklistJWT(ctx context.Context, jti string, ttlSeconds int) error {
	m.blacklist[jti] = true
	return nil
}

func (m *mockCacheRepository) IsJWTBlacklisted(ctx context.Context, jti string) (bool, error) {
	return m.blacklist[jti], nil
}

func (m *mockCacheRepository) IncrLoginAttempt(ctx context.Context, ip string, ttlSeconds int) (int, error) {
	if m.onIncrAttempts != nil {
		return m.onIncrAttempts(ip, ttlSeconds)
	}
	m.attempts[ip]++
	return m.attempts[ip], nil
}

func (m *mockCacheRepository) GetProfile(ctx context.Context, userID string) (string, error) {
	return m.profiles[userID], nil
}

func (m *mockCacheRepository) SetProfile(ctx context.Context, userID string, profileJSON string, ttlSeconds int) error {
	m.profiles[userID] = profileJSON
	return nil
}

func (m *mockCacheRepository) DeleteProfile(ctx context.Context, userID string) error {
	delete(m.profiles, userID)
	return nil
}
