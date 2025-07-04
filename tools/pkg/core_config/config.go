package coreconfig

import (
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
	"log"
)

// App параметры для запуска приложения
type App struct {
	Addr  string `envconfig:"APP_ADDR" required:"false" default:"0.0.0.0:9000"` // URL of the application
	Debug bool   `envconfig:"APP_DEBUG" default:"false"`
}

// Logging конфиг для создания логгера
type Logging struct {
	Level string `envconfig:"LOG_LEVEL" default:"debug"`
	File  string `envconfig:"LOG_FILE"`
	DSN   string `envconfig:"LOG_DSN"`
}

// Database конфигурация подключения к БД
type Database struct {
	URI      string `envconfig:"DB_URI"`
	Host     string `envconfig:"DB_HOST"`
	Port     int    `envconfig:"DB_PORT"`
	User     string `envconfig:"DB_USER"`
	Password string `envconfig:"DB_PASSWORD"`
	Name     string `envconfig:"DB_NAME"`
}

// Load метод для чтения конфига из окружения или .env файла
func Load(cfg interface{}, envNamespace string) error {
	// Load environment variables from the .env file
	if err := godotenv.Load(".env"); err != nil {
		log.Println("config file is not exists")
	}

	// Parse environment variables into the Config struct
	if err := envconfig.Process(envNamespace, cfg); err != nil {
		log.Fatalf("config not loaded: %s", err)
		return nil
	}

	// Return the loaded configuration
	return nil
}
