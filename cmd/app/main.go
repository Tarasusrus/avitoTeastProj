package main

import (
	"avitoTeastProj/internal/db"
	"avitoTeastProj/internal/models"
	"avitoTeastProj/internal/server"
	"go.uber.org/zap"
	"log"
)

//TODO: config()

// TODO: logger(slog)

// TODO: установка и настройка докера

// TODO:Установка и настройка БД

// TODO:Инициализация репозитория

// TODO:создание сущностей(User, Account, Product...?)

// TODO: структура БД (gorm)

// TODO: Создание методов реализующих бизнес логику

// TODO: Api. Создание ручек для методов

// TODO: Test. Разработка тестов(юнит, мод, нагр.)

// TODO: Opt-n. Кешировние, утечки памяти, (grafana)

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
		&models.Order{},
		&models.Reserve{},
		&models.ReportEntry{},
		&models.Service{})
	if err != nil {
		sugar.Fatalf("Migration failed: %v", err)
	}
	sugar.Info("Success migration")

	err = server.RunServer(data)
	if err != nil {
		sugar.Fatalf("Failed to run server: %v", err)
	}
}
