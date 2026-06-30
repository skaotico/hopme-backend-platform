package postgres

import (
	"context"
	"database/sql"
	"errors"
	"c4-auth/internal/domain/model"
	"c4-auth/internal/domain/port"
)

type userRepository struct {
	db *sql.DB
}

// NewUserRepository crea una instancia del adaptador de persistencia PostgreSQL
func NewUserRepository(db *sql.DB) port.UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, user *model.User) error {
	query := `
		INSERT INTO auth.users (username, email, password_hash, display_name, avatar_url, is_active, is_verified)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
		RETURNING id, created_at, updated_at
	`
	err := r.db.QueryRowContext(
		ctx, query,
		user.Username, user.Email, user.PasswordHash, user.DisplayName, user.AvatarURL, user.IsActive, user.IsVerified,
	).Scan(&user.ID, &user.CreatedAt, &user.UpdatedAt)

	if err != nil {
		return err
	}
	return nil
}

func (r *userRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, is_active, is_verified, last_login_at, created_at, updated_at
		FROM auth.users
		WHERE email = $1
	`
	user := &model.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, email).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName, &user.AvatarURL,
		&user.IsActive, &user.IsVerified, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}

	// Recuperar roles y sus permisos
	roles, err := r.fetchUserRolesAndPermissions(ctx, user.ID)
	if err == nil {
		user.Roles = roles
	}

	return user, nil
}

func (r *userRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, is_active, is_verified, last_login_at, created_at, updated_at
		FROM auth.users
		WHERE username = $1
	`
	user := &model.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, username).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName, &user.AvatarURL,
		&user.IsActive, &user.IsVerified, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}

	// Recuperar roles y sus permisos
	roles, err := r.fetchUserRolesAndPermissions(ctx, user.ID)
	if err == nil {
		user.Roles = roles
	}

	return user, nil
}

func (r *userRepository) FindByID(ctx context.Context, id string) (*model.User, error) {
	query := `
		SELECT id, username, email, password_hash, display_name, avatar_url, is_active, is_verified, last_login_at, created_at, updated_at
		FROM auth.users
		WHERE id = $1
	`
	user := &model.User{}
	var lastLogin sql.NullTime

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&user.ID, &user.Username, &user.Email, &user.PasswordHash, &user.DisplayName, &user.AvatarURL,
		&user.IsActive, &user.IsVerified, &lastLogin, &user.CreatedAt, &user.UpdatedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrUserNotFound
		}
		return nil, err
	}

	if lastLogin.Valid {
		user.LastLoginAt = &lastLogin.Time
	}

	// Recuperar roles y sus permisos
	roles, err := r.fetchUserRolesAndPermissions(ctx, user.ID)
	if err == nil {
		user.Roles = roles
	}

	return user, nil
}

func (r *userRepository) AssignRole(ctx context.Context, userID, roleCode string) error {
	query := `
		INSERT INTO iam.user_roles (user_id, role_id)
		VALUES ($1, (SELECT id FROM iam.roles WHERE code = $2))
		ON CONFLICT (user_id, role_id) DO NOTHING
	`
	_, err := r.db.ExecContext(ctx, query, userID, roleCode)
	return err
}

func (r *userRepository) FindRoleByCode(ctx context.Context, code string) (*model.Role, error) {
	query := `
		SELECT id, code, name, description, is_system
		FROM iam.roles
		WHERE code = $1
	`
	role := &model.Role{}
	err := r.db.QueryRowContext(ctx, query, code).Scan(
		&role.ID, &role.Code, &role.Name, &role.Description, &role.IsSystem,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, model.ErrRoleNotFound
		}
		return nil, err
	}

	// Recuperar los permisos asociados al rol
	perms, err := r.fetchRolePermissions(ctx, role.ID)
	if err == nil {
		role.Permissions = perms
	}

	return role, nil
}

// fetchUserRolesAndPermissions recupera los roles del usuario e inyecta sus permisos granulares
func (r *userRepository) fetchUserRolesAndPermissions(ctx context.Context, userID string) ([]model.Role, error) {
	query := `
		SELECT 
			r.id, r.code, r.name, r.description, r.is_system,
			p.id, p.code, p.description
		FROM iam.user_roles ur
		INNER JOIN iam.roles r ON r.id = ur.role_id
		LEFT JOIN iam.role_permissions rp ON rp.role_id = r.id
		LEFT JOIN iam.permissions p ON p.id = rp.permission_id
		WHERE ur.user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	rolesMap := make(map[string]*model.Role)

	for rows.Next() {
		var rID, rCode, rName, rDesc string
		var rIsSystem bool
		var pID, pCode, pDesc sql.NullString

		err := rows.Scan(&rID, &rCode, &rName, &rDesc, &rIsSystem, &pID, &pCode, &pDesc)
		if err != nil {
			return nil, err
		}

		role, exists := rolesMap[rID]
		if !exists {
			role = &model.Role{
				ID:          rID,
				Code:        rCode,
				Name:        rName,
				Description: rDesc,
				IsSystem:    rIsSystem,
				Permissions: []model.Permission{},
			}
			rolesMap[rID] = role
		}

		if pID.Valid {
			role.Permissions = append(role.Permissions, model.Permission{
				ID:          pID.String,
				Code:        pCode.String,
				Description: pDesc.String,
			})
		}
	}

	roles := make([]model.Role, 0, len(rolesMap))
	for _, role := range rolesMap {
		roles = append(roles, *role)
	}

	return roles, nil
}

// fetchRolePermissions busca los permisos de un rol específico
func (r *userRepository) fetchRolePermissions(ctx context.Context, roleID string) ([]model.Permission, error) {
	query := `
		SELECT p.id, p.code, p.description
		FROM iam.permissions p
		INNER JOIN iam.role_permissions rp ON rp.permission_id = p.id
		WHERE rp.role_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, roleID)
	if err != nil {
		return nil, err
	}
	defer func() { _ = rows.Close() }()

	var perms []model.Permission
	for rows.Next() {
		var perm model.Permission
		err := rows.Scan(&perm.ID, &perm.Code, &perm.Description)
		if err != nil {
			return nil, err
		}
		perms = append(perms, perm)
	}

	return perms, nil
}
