-- Habilitar extensión pgcrypto para gen_random_uuid()
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Creación de Esquemas
CREATE SCHEMA IF NOT EXISTS auth;
CREATE SCHEMA IF NOT EXISTS iam;
CREATE SCHEMA IF NOT EXISTS core;

-- ==========================================
-- ESQUEMA: AUTH
-- ==========================================

-- Tabla auth.users
CREATE TABLE IF NOT EXISTS auth.users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    username VARCHAR(50) NOT NULL UNIQUE,
    email VARCHAR(255) NOT NULL UNIQUE,
    password_hash TEXT NOT NULL,
    display_name VARCHAR(100),
    avatar_url TEXT,
    is_active BOOLEAN NOT NULL DEFAULT TRUE,
    is_verified BOOLEAN NOT NULL DEFAULT FALSE,
    last_login_at TIMESTAMP,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- ==========================================
-- ESQUEMA: IAM
-- ==========================================

-- Tabla iam.roles
CREATE TABLE IF NOT EXISTS iam.roles (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_system BOOLEAN NOT NULL DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabla iam.permissions
CREATE TABLE IF NOT EXISTS iam.permissions (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(100) NOT NULL UNIQUE,
    description TEXT,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabla iam.user_roles (Relación Muchos a Muchos Users <-> Roles)
CREATE TABLE IF NOT EXISTS iam.user_roles (
    user_id UUID NOT NULL,
    role_id UUID NOT NULL,
    assigned_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, role_id),
    CONSTRAINT fk_user_roles_user
        FOREIGN KEY(user_id)
        REFERENCES auth.users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_roles_role
        FOREIGN KEY(role_id)
        REFERENCES iam.roles(id)
        ON DELETE CASCADE
);

-- Tabla iam.role_permissions (Relación Muchos a Muchos Roles <-> Permissions)
CREATE TABLE IF NOT EXISTS iam.role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    PRIMARY KEY(role_id, permission_id),
    CONSTRAINT fk_role_permissions_role
        FOREIGN KEY(role_id)
        REFERENCES iam.roles(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_role_permissions_permission
        FOREIGN KEY(permission_id)
        REFERENCES iam.permissions(id)
        ON DELETE CASCADE
);

-- ==========================================
-- ESQUEMA: CORE
-- ==========================================

-- Tabla core.modules
CREATE TABLE IF NOT EXISTS core.modules (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    code VARCHAR(50) NOT NULL UNIQUE,
    name VARCHAR(100) NOT NULL,
    description TEXT,
    is_enabled BOOLEAN NOT NULL DEFAULT TRUE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW()
);

-- Tabla core.user_modules (Relación Muchos a Muchos Users <-> Modules)
CREATE TABLE IF NOT EXISTS core.user_modules (
    user_id UUID NOT NULL,
    module_id UUID NOT NULL,
    enabled_at TIMESTAMP NOT NULL DEFAULT NOW(),
    PRIMARY KEY(user_id, module_id),
    CONSTRAINT fk_user_modules_user
        FOREIGN KEY(user_id)
        REFERENCES auth.users(id)
        ON DELETE CASCADE,
    CONSTRAINT fk_user_modules_module
        FOREIGN KEY(module_id)
        REFERENCES core.modules(id)
        ON DELETE CASCADE
);

-- ==========================================
-- ÍNDICES RECOMENDADOS
-- ==========================================
CREATE INDEX IF NOT EXISTS idx_users_email ON auth.users(email);
CREATE INDEX IF NOT EXISTS idx_users_username ON auth.users(username);
CREATE INDEX IF NOT EXISTS idx_roles_code ON iam.roles(code);
CREATE INDEX IF NOT EXISTS idx_permissions_code ON iam.permissions(code);
