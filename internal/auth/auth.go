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
	UserId string `json:"username"`
	jwt.RegisteredClaims
}

type UserToken struct {
	Token string
}

var JwtKey string

func InitjwtKey(auth *env_config.Auth) {
	JwtKey = auth.JwtKey
}

func GetJwtKey() string {
	return JwtKey
}
func CreateToken(id string) (string, error) {
	var jwtKey = []byte(JwtKey)

	claims := CustomClaims{
		// subject username
		id,
		jwt.RegisteredClaims{
			// expiry time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	// Creates a new token with the specified signing method and claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	//Creates and returns a complete, signed JWT
	ss, err := token.SignedString(jwtKey)
	if err != nil {
		return "", errors.New("Error signing token")
	}
	return ss, err
}
