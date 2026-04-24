package config

import (
	"os"
)

type Config struct {
	Port               string
	Debug              string
	JWTSecret          string
	RefreshTokenSecret string
	IdentityProviderID string
	FrontendURL 	   string
	Database           DatabaseConfig
	Redis              RedisConfig
	Email              EmailConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Username string
	Password string
	Name     string
	SSLMode  string
}

type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

type EmailConfig struct {
	EmailUsername string
	EmailPassword string
	SmtpHost      string
	SmtpPort      string
}

func Load() *Config {
	return &Config{
		Port:               GetEnv("PORT", "8080"),
		Debug:              GetEnv("DEBUG", "false"),	
		JWTSecret:          GetEnv("JWT_SECRET", "secret"),
		RefreshTokenSecret: GetEnv("REFRESH_TOKEN_SECRET", "refresh-secret"),
		FrontendURL: 		GetEnv("FRONTEND_URL", "http://localhost:3000"),
		IdentityProviderID: GetEnv("IDENTITY_PROVIDER_ID", "1"),
		Database: DatabaseConfig{
			Host:     GetEnv("DB_HOST", "localhost"),
			Port:     GetEnv("DB_PORT", "3306"),
			Username: GetEnv("DB_USER", "root"),
			Password: GetEnv("DB_PASS", ""),
			Name:     GetEnv("DB_NAME", ""),
			SSLMode:  GetEnv("DB_SSLMODE", "disable"),
		},
		Redis: RedisConfig{
			Addr:     GetEnv("REDIS_ADDR", "localhost") + ":" + GetEnv("REDIS_PORT", "6379"),
			Password: GetEnv("REDIS_PASSWORD", ""),
			DB:       0,
		},
		Email: EmailConfig{
			EmailUsername: GetEnv("EMAIL_USERNAME", ""),
			EmailPassword: GetEnv("EMAIL_PASSWORD", ""),
			SmtpHost:     GetEnv("EMAIL_SMTP_HOST", ""),
			SmtpPort:     GetEnv("EMAIL_SMTP_PORT", ""),
		},
	}
}

func GetEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}