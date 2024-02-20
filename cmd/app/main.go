package main

import (
	"avitoTeastProj/internal/db"
	"avitoTeastProj/internal/models"
	"avitoTeastProj/internal/server"
	"go.uber.org/zap"
	"log"
)

func main() {
	logger, err := zap.NewDevelopment()
	if err != nil {
		log.Fatalf("Can't initialize zap logger: %v", err)
	}

	defer logger.Sync() // Обязательно для корректного завершения работы логгера

	sugar := logger.Sugar()

	sugar.Infow("Zap Logger initialized successfully",
		"version", "v1.0.0", "mode", "production",
	)

	cfg, err := db.LoadConfig("./configs/config.yaml")
	if err != nil {
		sugar.Fatalf("Failed to load configs %v", err)
	}

	data, err := db.ConnectToDb(cfg)
	if err != nil {
		sugar.Fatalf("Failed to connect DB %v", err)
	}

	err = data.AutoMigrate(&models.User{},
		&models.Reserve{},
		&models.ReportEntry{})
	if err != nil {
		sugar.Fatalf("Migration failed: %v", err)
	}
	sugar.Info("Success migration")

	err = server.RunServer(data, sugar)
	if err != nil {
		sugar.Fatalf("Failed to run server: %v", err)
	}
}
