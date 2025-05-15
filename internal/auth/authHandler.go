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
			gin.H{
				"error": "Token is not found",
			},
		)
		c.Abort()
		return
	}

	c.JSON(
		http.StatusOK,
		gin.H{
			"OK": token,
		},
	)
}
