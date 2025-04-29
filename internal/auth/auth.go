package auth

import (
	"fmt"
	"jennyfood/models"
	"math/rand"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func CryptPassword(password string) []byte {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), rand.Int())
	if err != nil {
		fmt.Printf("%v", err)
		return []byte(password)
	}
	return cryptedPassword
}

func GenerateToken(userdata models.UserData) string {
	payload := jwt.MapClaims{
		"sub": userdata.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signed, err := token.SignedString(os.Getenv("SECRET"))
	if err != nil {
		fmt.Printf("%v", err)
		return ""
	}
	return signed
}
