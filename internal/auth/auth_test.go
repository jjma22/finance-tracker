package auth_test

import (
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/jjma22/finance-tracker/internal/auth"
	env_config "github.com/jjma22/finance-tracker/internal/config"
)

func TestCreateToken(t *testing.T) {
	user := "James"

	testAuth := env_config.Auth{
		JwtKey: "abcdefg12345",
	}

	auth.InitjwtKey(&testAuth)
	token, err := auth.CreateToken(user)

	if err != nil {
		t.Fatalf("Errur crearing jwt for user: %v", err)
	}

	claims := auth.CustomClaims{
		// subject username
		user,
		jwt.RegisteredClaims{
			// expiry time
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}

	parsedToken, err := jwt.ParseWithClaims(token, &claims, func(token *jwt.Token) (any, error) {
		return []byte("abcdefg12345"), nil
	})

	if err != nil {
		t.Fatalf("Error parsing token: %v", err)
	}

	if parsedToken.Valid != true {
		t.Fatalf("jwt token is not valid: %v", err)
	}

}
