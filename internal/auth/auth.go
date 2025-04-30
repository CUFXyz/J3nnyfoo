package auth

import (
	"fmt"
	"jennyfood/models"
	"log"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jmoiron/sqlx"
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

func ClaimToken(us models.RegisterData, db *sqlx.DB) string {
	newus := models.UserData{}
	newus.Token = GenerateToken(us)

	_, err := db.Exec("UPDATE users SET token = $1 WHERE email = $2;", newus.Token, us.Email)
	if err != nil {
		log.Fatalf("Error due executing updating row")
		return ""
	}
	return newus.Token
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
