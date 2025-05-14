package auth

import (
	"fmt"
	"jennyfood/internal/config"
	"jennyfood/models"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthInstance struct {
	AuthCfg config.AuthConfig
}

// Dehashing password and comparing them to pass user or not
func (ai *AuthInstance) AuthFunc(hashedPass []byte, origPass []byte) error {
	if err := bcrypt.CompareHashAndPassword(hashedPass, origPass); err != nil {
		return err
	}

	return nil
}

// Hashing password for some sort of security
func (ai *AuthInstance) CryptPassword(password string) []byte {
	cryptedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		fmt.Printf("Error due crypting password, returning noncrypted")
		return []byte(password)
	}
	return cryptedPassword
}

// Generating token/cookie for user
func (ai *AuthInstance) GenerateToken(us models.RegisterData) string {
	payload := jwt.MapClaims{
		"sub": us.Email,
		"exp": time.Now().Add(time.Hour * 24).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	signed, err := token.SignedString([]byte(ai.AuthCfg.Secret))
	if err != nil {
		fmt.Printf("Error due signing string")
		return ""
	}
	return signed
}
