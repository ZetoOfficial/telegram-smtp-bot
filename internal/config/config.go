package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	TelegramToken string
	SMTPHost      string
	SMTPPort      int
	SMTPUsername  string
	SMTPPassword  string
}

func LoadConfig(path string) (*Config, error) {
	err := godotenv.Load(path)
	if err != nil {
		log.Printf("Warning: could not load config from %s, using environment variables", path)
	}

	port := getEnv("SMTP_PORT", "587")
	intPort, err := strconv.Atoi(port)
	if err != nil {
		return nil, fmt.Errorf("invalid SMTP_PORT: %w", err)
	}
	config := &Config{
		TelegramToken: getEnv("TELEGRAM_TOKEN", ""),
		SMTPHost:      getEnv("SMTP_HOST", "smtp.yandex.com"),
		SMTPPort:      intPort,
		SMTPUsername:  getEnv("SMTP_USERNAME", ""),
		SMTPPassword:  getEnv("SMTP_PASSWORD", ""),
	}

	if config.TelegramToken == "" ||
		config.SMTPHost == "" ||
		config.SMTPPort == 0 ||
		config.SMTPUsername == "" ||
		config.SMTPPassword == "" {
		return nil, fmt.Errorf("env is empty")
	}

	return config, nil
}

func getEnv(key, defaultValue string) string {
	value, exists := os.LookupEnv(key)
	if !exists {
		return defaultValue
	}
	return value
}
