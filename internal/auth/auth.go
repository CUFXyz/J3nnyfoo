package auth

import (
	"fmt"
	"jennyfood/models"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

func AuthFunc(hashedPass []byte, origPass []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPass, origPass); err != nil {
		return err
	}

	return nil
}

func CryptPassword(password string) []byte {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error due crypting password, returning noncrypted")
		return []byte(password)
	}
	return cryptedPassword
}

func GenerateToken(us models.RegisterData) string {
	payload := jwt.MapClaims{
		"sub": us.Email,
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
