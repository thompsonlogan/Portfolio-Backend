package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	AppEnv      string
	Port        string
	DBHost      string
	DBUser      string
	DBPassword  string
	DBName      string
	DBPort      string
	DBSSLMode   string
	FrontendURL string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	cfg := &Config{
		AppEnv:      getRequiredEnv("APP_ENV"),
		Port:        getRequiredEnv("PORT"),
		DBHost:      getRequiredEnv("DB_HOST"),
		DBUser:      getRequiredEnv("DB_USER"),
		DBPassword:  getRequiredEnv("DB_PASSWORD"),
		DBName:      getRequiredEnv("DB_NAME"),
		DBPort:      getRequiredEnv("DB_PORT"),
		DBSSLMode:   getRequiredEnv("DB_SSLMODE"),
		FrontendURL: getRequiredEnv("FRONTEND_URL"),
	}

	return cfg
}

func getRequiredEnv(key string) string {
	val, exists := os.LookupEnv(key)
	if !exists || val == "" {
		log.Fatalf("environment variable %s is required but not set", key)
	}
	return val
}

func (c *Config) PostgresDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=UTC",
		c.DBHost, c.DBUser, c.DBPassword, c.DBName, c.DBPort, c.DBSSLMode,
	)
}

func (c *Config) IsDev() bool {
	return c.AppEnv != "production"
}

func (c *Config) Log() {
	log.Printf("Environment: %s, DB: %s@%s:%s/%s, FrontendURL: %s",
		c.AppEnv, c.DBUser, c.DBHost, c.DBPort, c.DBName, c.FrontendURL)
}
