CREATE SCHEMA IF NOT EXISTS "catalogo";

CREATE TABLE IF NOT EXISTS "catalogo"."especie_arbol" (
  "id" uuid PRIMARY KEY,
  "nombre_comun" varchar(200),
  "nombre_cientifico" varchar(200) NOT NULL,
  "familia" varchar(200),
  "descripcion" text,
  "activo" boolean DEFAULT true
);

CREATE TABLE IF NOT EXISTS "catalogo"."estado_arbol" (
  "id" uuid PRIMARY KEY,
  "codigo" varchar(50) UNIQUE,
  "nombre" varchar(100),
  "descripcion" text
);

CREATE TABLE IF NOT EXISTS "catalogo"."estado_estanque" (
  "id" uuid PRIMARY KEY,
  "codigo" varchar(50) UNIQUE,
  "nombre" varchar(100)
);

CREATE TABLE IF NOT EXISTS "catalogo"."estado_agua" (
  "id" uuid PRIMARY KEY,
  "codigo" varchar(50) UNIQUE,
  "nombre" varchar(100)
);

CREATE TABLE IF NOT EXISTS "catalogo"."tipo_sensor" (
  "id" uuid PRIMARY KEY,
  "codigo" varchar(50) UNIQUE,
  "nombre" varchar(100),
  "unidad_medida" varchar(50),
  "descripcion" text
);
