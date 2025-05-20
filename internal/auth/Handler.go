package auth

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

func (ai *AuthInstance) AuthHandler(c *gin.Context) {
	token, exists := c.Cookie("token")
	var resultdb string
	if exists != nil {
		c.JSON(
			http.StatusForbidden,
			"Service is available for logged users",
		)
		c.Abort()
		return
	}
	result, err := ai.Cache.GetValue(token)
	fmt.Println()
	fmt.Println(result)
	fmt.Println()
	if err != nil {
		c.JSON(
			http.StatusBadGateway,
			gin.H{
				"result": result,
				"Error":  error.Error(err),
			},
		)
		c.Abort()
		return
	}
	queryresult, err := ai.Pgdb.Query("SELECT email FROM users WHERE email = $1;", result)
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			error.Error(err),
		)
		c.Abort()
		return
	}
	err = queryresult.Scan(&resultdb)
	fmt.Println()
	fmt.Println(resultdb)
	fmt.Println()
	if err != nil {
		c.JSON(
			http.StatusBadRequest,
			error.Error(err),
		)
		c.Abort()
		return
	}
	if result == resultdb {
		c.JSON(
			http.StatusOK,
			gin.H{
				"Auth": "OK",
			},
		)
	}

}
