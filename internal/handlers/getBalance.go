package handlers

import (
	"avitoTeastProj/internal/service"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

func GetBalanceHandler(c *gin.Context) {

	userSrvice, exists := c.MustGet("userService").(*service.UserService)
	if !exists {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "service not available"})
		return
	}

	idParam := c.Query("id")

	userid, err := strconv.ParseUint(idParam, 10, 64)
	if err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, err)
		return
	}

	balance, err := userSrvice.GetBalance(uint(userid))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error getting balance"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})

}
