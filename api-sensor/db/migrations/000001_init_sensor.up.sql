-- Habilitar extensión pgcrypto si no está habilitada
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Creación de Esquema
CREATE SCHEMA IF NOT EXISTS "iot";

-- Tabla: iot.sensor
CREATE TABLE IF NOT EXISTS "iot"."sensor" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "codigo" VARCHAR(100) UNIQUE,
    "tipo_sensor_id" UUID NOT NULL,
    "arbol_id" UUID,
    "estanque_id" UUID,
    "fabricante" VARCHAR(200),
    "modelo" VARCHAR(200),
    "fecha_instalacion" TIMESTAMP,
    "activo" BOOLEAN DEFAULT true
);

-- Tabla: iot.lectura_sensor
CREATE TABLE IF NOT EXISTS "iot"."lectura_sensor" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "sensor_id" UUID NOT NULL,
    "valor" NUMERIC(18, 6) NOT NULL,
    "fecha_lectura" TIMESTAMP NOT NULL DEFAULT NOW(),
    "bateria_porcentaje" INTEGER,
    "observacion" TEXT,
    CONSTRAINT fk_lectura_sensor
        FOREIGN KEY (sensor_id)
        REFERENCES "iot"."sensor"(id)
        ON DELETE CASCADE
);

-- Tabla: iot.alerta
CREATE TABLE IF NOT EXISTS "iot"."alerta" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "sensor_id" UUID NOT NULL,
    "tipo" VARCHAR(100) NOT NULL,
    "valor_detectado" NUMERIC(18, 6),
    "umbral" NUMERIC(18, 6),
    "fecha_alerta" TIMESTAMP NOT NULL DEFAULT NOW(),
    "estado" VARCHAR(50) NOT NULL DEFAULT 'activa',
    "observacion" TEXT,
    CONSTRAINT fk_alerta_sensor
        FOREIGN KEY (sensor_id)
        REFERENCES "iot"."sensor"(id)
        ON DELETE CASCADE
);
