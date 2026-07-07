-- Habilitar extensión pgcrypto si no está habilitada
CREATE EXTENSION IF NOT EXISTS "pgcrypto";

-- Creación de Esquema
CREATE SCHEMA IF NOT EXISTS "flora";

-- Tabla: flora.arbol
CREATE TABLE IF NOT EXISTS "flora"."arbol" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "zona_id" UUID NOT NULL,
    "especie_id" UUID,
    "estado_id" UUID,
    "codigo" VARCHAR(100),
    "latitud" DOUBLE PRECISION,
    "longitud" DOUBLE PRECISION,
    "edad_estimada_anios" INTEGER,
    "altura_m" DOUBLE PRECISION,
    "ancho_copa_m" DOUBLE PRECISION,
    "diametro_tronco_cm" DOUBLE PRECISION,
    "fecha_plantacion" DATE,
    "fecha_registro" TIMESTAMP,
    "observaciones" TEXT,
    "fecha_creacion" TIMESTAMP NOT NULL DEFAULT NOW(),
    "fecha_actualizacion" TIMESTAMP
);

-- Tabla: flora.historial_estado_arbol
CREATE TABLE IF NOT EXISTS "flora"."historial_estado_arbol" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "arbol_id" UUID NOT NULL,
    "estado_id" UUID NOT NULL,
    "observacion" TEXT,
    "fecha_registro" TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_historial_arbol
        FOREIGN KEY (arbol_id)
        REFERENCES "flora"."arbol"(id)
        ON DELETE CASCADE
);

-- Tabla: flora.medicion_arbol
CREATE TABLE IF NOT EXISTS "flora"."medicion_arbol" (
    "id" UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    "arbol_id" UUID NOT NULL,
    "altura_m" DOUBLE PRECISION,
    "ancho_copa_m" DOUBLE PRECISION,
    "diametro_tronco_cm" DOUBLE PRECISION,
    "observaciones" TEXT,
    "fecha_medicion" TIMESTAMP NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_medicion_arbol
        FOREIGN KEY (arbol_id)
        REFERENCES "flora"."arbol"(id)
        ON DELETE CASCADE
);
