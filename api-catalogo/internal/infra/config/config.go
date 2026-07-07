package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
)

// AppConfig almacena toda la configuración leída de variables de entorno
type AppConfig struct {
	Port string

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RunMigrations bool

	LogLevel    string
	Environment string
}

// loadDotEnv lee un archivo .env local y carga sus variables de entorno
func loadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

// Load lee las variables de entorno e inicializa la configuración de la aplicación
func Load() *AppConfig {
	loadDotEnv()

	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		dbPort = 5432
	}

	runMigrationsStr := getEnv("RUN_MIGRATIONS", "false")
	runMigrations, err := strconv.ParseBool(runMigrationsStr)
	if err != nil {
		runMigrations = false
	}

	return &AppConfig{
		Port: getEnv("PORT", "9094"),

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "adminpassword"),
		DBName:     getEnv("DB_DATABASE", getEnv("DB_NAME", "desaSkaotico")),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RunMigrations: runMigrations,

		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "production"),
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
