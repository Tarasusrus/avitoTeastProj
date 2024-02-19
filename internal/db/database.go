package db

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
)

type Config struct {
	Database struct {
		Host     string `env:"DB_HOST" env-default:"127.0.0.1"`
		User     string `env:"DB_USER" env-default:"your_database_user"`
		Password string `env:"DB_PASSWORD" env-default:"your_database_password"`
		Name     string `env:"DB_NAME" env-default:"your_database_name"`
		Port     int    `env:"DB_PORT" env-default:"3306"`
	} `env:"DATABASE"`
}

// LoadConfig загружает конфигурацию из файла
func LoadConfig(filePath string) (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig(filePath, &cfg); err != nil {
		log.Fatalf("Error reading config file: %s", err)
		return nil, err
	}
	return &cfg, nil

}

func ConnectToDb(cfg *Config) (*gorm.DB, error) {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.Name,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Err in ConnectToDb: %s", err)
		return nil, err
	}
	return db, err
}
