-- Borrar tablas que poseen claves foráneas primero
DROP TABLE IF EXISTS core.user_modules CASCADE;
DROP TABLE IF EXISTS iam.role_permissions CASCADE;
DROP TABLE IF EXISTS iam.user_roles CASCADE;

-- Borrar tablas base
DROP TABLE IF EXISTS core.modules CASCADE;
DROP TABLE IF EXISTS iam.permissions CASCADE;
DROP TABLE IF EXISTS iam.roles CASCADE;
DROP TABLE IF EXISTS auth.users CASCADE;

-- Borrar esquemas
DROP SCHEMA IF EXISTS core CASCADE;
DROP SCHEMA IF EXISTS iam CASCADE;
DROP SCHEMA IF EXISTS auth CASCADE;
