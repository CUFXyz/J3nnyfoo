package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ai *AuthInstance) AuthHandler(c *gin.Context) {
	token, exists := c.Cookie("token")
	if exists != nil {
		c.JSON(
			http.StatusForbidden,
			"Service is available for logged users",
		)
		c.Abort()
		return
	}

	result, err := ai.Cache.GetValue(token)
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			"Something wrong with app cache",
		)
		c.Abort()
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"token": result,
			"Auth":  "OK",
		},
	)
}
