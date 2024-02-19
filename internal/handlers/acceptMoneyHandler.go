package handlers

import (
	"avitoTeastProj/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
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
	userSrvice, exists := c.MustGet("userService").(*service.UserService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service not available"})
		return
	}

	bodyBytes, err := io.ReadAll(c.Request.Body)
	defer c.Request.Body.Close()
	if err != nil {
		http.Error(c.Writer, "Bad Request in bodyBytes", http.StatusBadRequest)
		return
	}
	var input ReserveMoneyInput
	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		http.Error(c.Writer, "Bad Request in Unmarshal", http.StatusBadRequest)
		return
	}

	err = userSrvice.AcceptMoney(input.Id, input.ServiceId, input.OrderId, input.Price)
	if err != nil {
		http.Error(c.Writer, "Internal Server Error in userService.AcceptMoney", http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
