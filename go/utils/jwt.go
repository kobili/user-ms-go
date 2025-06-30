package utils

import (
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"kobili/user-ms/entities"
)

func CreateJWTForUser(user entities.UserEntity) (string, error) {
	key := os.Getenv("JWT_SIGNING_KEY")
	if key == "" {
		return "", fmt.Errorf("JWT_SIGNING_KEY not set")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256,
		jwt.MapClaims{
			"iss": "kobili/user-ms",
			"sub": user.Id,
			"iat": time.Now().Unix(),
			"exp": time.Now().Add(15 * time.Minute).Unix(),
		},
	)

	signedToken, err := token.SignedString([]byte(key))
	if err != nil {
		return "", nil
	}

	return signedToken, nil
}
