package handlers

import (
	"avitoTeastProj/internal/service"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
)

type BalanceInput struct {
	Id      uint `json:"id"`
	Balance uint `json:"balance"`
}

func SetBalance(c *gin.Context) {

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
	var input BalanceInput

	err = json.Unmarshal(bodyBytes, &input)
	if err != nil {
		http.Error(c.Writer, "Bad Request in Unmarshal", http.StatusBadRequest)
		return
	}

	err = userSrvice.UpdateBalance(input.Id, input.Balance)
	if err != nil {
		http.Error(c.Writer, "Internal Server Error in userService.UpdateBalance", http.StatusInternalServerError)
		return
	}
	c.Writer.WriteHeader(http.StatusOK)
}
