package auth

import (
	"fmt"
	"jennyfood/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AuthFunc(hashedPass []byte, origPass []byte) (bool, error) {
	result := bcrypt.CompareHashAndPassword(hashedPass, origPass)
	if result != nil {
		return false, fmt.Errorf("password is not correct")
	}
	return true, nil
}

func CryptPassword(password string) []byte {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), 10)
	if err != nil {
		fmt.Printf("Error due crypting password, returning noncrypted")
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
	signed, err := token.SignedString([]byte(os.Getenv("SECRET")))
	if err != nil {
		fmt.Printf("Error due signing string")
		return ""
	}
	return signed
}
