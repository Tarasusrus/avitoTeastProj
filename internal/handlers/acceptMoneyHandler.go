package handlers

import (
	"avitoTeastProj/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io"
	"net/http"
)

type AcceptMoneyInput struct {
	Id        uint `json:"id"`
	Price     uint `json:"price"`
	ServiceId uint `json:"serviceId"`
	OrderId   uint `json:"orderId"`
}

func AcceptMoneyHandlers(c *gin.Context) {
	logger := c.MustGet("logger").(*zap.SugaredLogger)

	// Попытка извлечь экземпляр UserService из контекста
	userSrvice, exists := c.MustGet("userService").(*service.UserService)
	if !exists {
		logger.Error("UserService not available")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service not available"})
		return
	}

	// Чтение тела запроса
	bodyBytes, err := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		logger.Errorw("Error reading request body", "error", err)
		http.Error(c.Writer, "Bad Request in bodyBytes", http.StatusBadRequest)
		return
	}

	// Десериализация входных данных из тела запроса
	var input AcceptMoneyInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		logger.Errorw("Error unmarshaling request body", "error", err)
		http.Error(c.Writer, "Bad Request in Unmarshal", http.StatusBadRequest)
		return
	}

	// Выполнение операции принятия денежных средств
	err = userSrvice.AcceptMoney(input.Id, input.ServiceId, input.OrderId, input.Price)
	if err != nil {
		logger.Errorw("Error in userService.AcceptMoney", "error", err)
		http.Error(c.Writer, "Internal Server Error in userService.AcceptMoney", http.StatusInternalServerError)
		return
	}
	// Успешное завершение обработки запроса
	logger.Infow("Money accept successfully", "userId", input.Id, "price", input.Price)
	c.Writer.WriteHeader(http.StatusOK)
}
