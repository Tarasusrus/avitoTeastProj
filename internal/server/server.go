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

	r.Use(func(c *gin.Context) {
		c.Set("userService", userService)
		c.Next()
	})

	r.GET("/balance", handlers.GetBalanceHandler)
	r.POST("/setbalance", func(c *gin.Context) {
		handlers.SetBalance(c)
	})
	r.POST("/setbalance/reservmoney", func(c *gin.Context) {
		handlers.ReserveMoneyHandlers(c)
	})
	r.POST("/setbalance/acceptmoney", func(c *gin.Context) {
		handlers.AcceptMoneyHandlers(c)
	})

	//// Пример использования замыкания для передачи userService
	//r.POST("/userBalance", func(c *gin.Context) {
	//	// Извлечение userService из контекста запроса и передача его в хендлер
	//	userService, _ := c.MustGet("userService").(*service.UserService)
	//	handlers.SetBalance(c, userService)
	//})

	// Запуск сервера
	if err := r.Run(":8080"); err != nil {
		log.Fatal("Server Run Failed", err)
	}
	return nil
}
