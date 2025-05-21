package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ai *AuthInstance) AuthHandler(ctx *gin.Context) {
	token := ctx.GetHeader("Authorization")
	if token == "" {
		ctx.JSON(
			http.StatusBadRequest,
			gin.H{
				"error": "No token provided",
			},
		)
		return
	}

	_, err := ai.Cache.GetValue(token)
	if err != nil {
		ctx.JSON(
			http.StatusBadGateway,
			fmt.Sprintf("Something wrong with app cache: %v", err),
		)
		ctx.Abort()
		return
	}

	ctx.Next() // Pass it to the next handler.
}
