package main

import (
	"avitoTeastProj/internal/db"
	"avitoTeastProj/internal/models"
	"avitoTeastProj/internal/server"
	"log"
	"log/slog"
	"os"
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
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))

	cfg, err := db.LoadConfig("./configs/config.yaml")
	if err != nil {
		logger.Error("err in load config", err)
	}

	data, err := db.ConnectToDb(cfg)
	if err != nil {
		log.Fatal("err in ConnectToDb", err)
	}

	err = data.AutoMigrate(&models.User{},
		&models.Order{},
		&models.Reserve{},
		&models.ReportEntry{},
		&models.Service{})
	if err != nil {
		logger.Error("Migration failed: ", err)
	}
	logger.Info("Sucsess migratione")

	err = server.RunServer(data)
	if err != nil {
		log.Fatal("Failed to run server: ", err)
	}
}
