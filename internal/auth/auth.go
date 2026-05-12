package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	env_config "github.com/jjma22/finance-tracker/internal/config"
)

type User struct {
	Id       string
	Username string
	Password string
}

type CustomClaims struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

var JwtKey string

func InitjwtKey(auth *env_config.Auth) {
	JwtKey = auth.JwtKey
}

func createToken(u string) (string, error) {
	var jwtKey = []byte(JwtKey)

	claims := CustomClaims{
		u,
		jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("Error signing token")
	}
	return ss, err
}
