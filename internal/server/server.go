package server

import (
	"avitoTeastProj/internal/handlers"
	"avitoTeastProj/internal/service"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func RunServer(db *gorm.DB) error {
	userService := service.NewUserService(db)

	r := gin.Default()

	// Глобальный middleware для добавления userService в контекст каждого запроса
	r.Use(func(c *gin.Context) {
		c.Set("userService", userService)
		c.Next()
	})

	// Группировка ручек, связанных с балансом
	balanceRouters := r.Group("/balance")
	{
		balanceRouters.GET("/", handlers.GetBalanceHandler)
		balanceRouters.POST("/set", handlers.SetBalance)
	}

	// Группировка ручек, связанных с операциями над балансом
	operationRouters := r.Group("/operations")
	{
		operationRouters.POST("/reservemoney", handlers.ReserveMoneyHandlers)
		operationRouters.POST("/acceptmoney", handlers.AcceptMoneyHandlers)
	}

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed", err)
	}
	return nil
}
