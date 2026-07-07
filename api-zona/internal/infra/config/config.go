package config

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)

// AppConfig almacena toda la configuración leída de variables de entorno
type AppConfig struct {
	Port          string
	JWTSecret     string
	JWTExpiration time.Duration

	DBHost     string
	DBPort     int
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string

	RedisHost     string
	RedisPort     int
	RedisPassword string
	RedisDB       int

	TimeoutLogin    time.Duration
	TimeoutRegister time.Duration
	TimeoutMe       time.Duration

	RTExpiration time.Duration

	LoginMaxAttempts  int
	LoginBlockMinutes int

	RunMigrations bool

	// Observabilidad
	// LogLevel controla el nivel mínimo de logs emitidos ("debug", "info", "warn", "error").
	// En desarrollo usar "debug" para ver todos los pasos internos.
	LogLevel string
	// Environment identifica el entorno de ejecución ("development", "staging", "production").
	// En "development" los logs incluyen caller (archivo:línea) para facilitar depuración.
	Environment string

	// Configuración de Cookies
	CookieSecure bool
}

// loadDotEnv lee un archivo .env local y carga sus variables de entorno
func loadDotEnv() {
	file, err := os.Open(".env")
	if err != nil {
		return // Si no existe el archivo .env, no hacemos nada y continuamos
	}
	defer func() { _ = file.Close() }()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		// Ignorar líneas vacías o comentarios
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Quitar comillas simples o dobles
		if len(value) >= 2 && ((value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'')) {
			value = value[1 : len(value)-1]
		}

		// Registrar en el entorno si no está ya definida por el sistema
		if _, exists := os.LookupEnv(key); !exists {
			_ = os.Setenv(key, value)
		}
	}
}

// Load lee las variables de entorno e inicializa la configuración de la aplicación
func Load() *AppConfig {
	// Intentar cargar variables de entorno del archivo .env primero
	loadDotEnv()

	port := getEnv("PORT", "8080")
	jwtSecret := getEnv("JWT_SECRET", "1Dos3M0M13")
	jwtExpStr := getEnv("JWT_EXPIRATION_HOURS", "24")

	jwtExpHours, err := strconv.Atoi(jwtExpStr)
	if err != nil {
		jwtExpHours = 24
	}

	dbPortStr := getEnv("DB_PORT", "5432")
	dbPort, err := strconv.Atoi(dbPortStr)
	if err != nil {
		dbPort = 5432
	}

	timeoutLoginStr := getEnv("TIMEOUT_LOGIN", "15")
	timeoutLogin, err := strconv.Atoi(timeoutLoginStr)
	if err != nil {
		timeoutLogin = 15
	}

	timeoutRegisterStr := getEnv("TIMEOUT_REGISTER", "30")
	timeoutRegister, err := strconv.Atoi(timeoutRegisterStr)
	if err != nil {
		timeoutRegister = 30
	}

	timeoutMeStr := getEnv("TIMEOUT_ME", "10")
	timeoutMe, err := strconv.Atoi(timeoutMeStr)
	if err != nil {
		timeoutMe = 10
	}

	runMigrationsStr := getEnv("RUN_MIGRATIONS", "false")
	runMigrations, err := strconv.ParseBool(runMigrationsStr)
	if err != nil {
		runMigrations = false
	}

	rtExpStr := getEnv("RT_EXPIRATION_DAYS", "30")
	rtExpDays, err := strconv.Atoi(rtExpStr)
	if err != nil {
		rtExpDays = 30
	}

	loginMaxAttemptsStr := getEnv("LOGIN_MAX_ATTEMPTS", "5")
	loginMaxAttempts, err := strconv.Atoi(loginMaxAttemptsStr)
	if err != nil {
		loginMaxAttempts = 5
	}

	loginBlockMinutesStr := getEnv("LOGIN_BLOCK_MINUTES", "15")
	loginBlockMinutes, err := strconv.Atoi(loginBlockMinutesStr)
	if err != nil {
		loginBlockMinutes = 15
	}

	redisPortStr := getEnv("REDIS_PORT", "6379")
	redisPort, err := strconv.Atoi(redisPortStr)
	if err != nil {
		redisPort = 6379
	}

	redisDBStr := getEnv("REDIS_DB", "0")
	redisDB, err := strconv.Atoi(redisDBStr)
	if err != nil {
		redisDB = 0
	}

	return &AppConfig{
		Port:          port,
		JWTSecret:     jwtSecret,
		JWTExpiration: time.Duration(jwtExpHours) * time.Hour,
		RTExpiration:  time.Duration(rtExpDays) * 24 * time.Hour,

		LoginMaxAttempts:  loginMaxAttempts,
		LoginBlockMinutes: loginBlockMinutes,

		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     dbPort,
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "adminpassword"),
		DBName:     getEnv("DB_DATABASE", getEnv("DB_NAME", "desaSkaotico")),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),

		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     redisPort,
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,

		TimeoutLogin:    time.Duration(timeoutLogin) * time.Second,
		TimeoutRegister: time.Duration(timeoutRegister) * time.Second,
		TimeoutMe:       time.Duration(timeoutMe) * time.Second,

		RunMigrations: runMigrations,

		// Observabilidad: por defecto info en producción, configurable vía env
		LogLevel:    getEnv("LOG_LEVEL", "info"),
		Environment: getEnv("ENVIRONMENT", "production"),
		CookieSecure: getEnv("ENVIRONMENT", "production") != "development",
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
