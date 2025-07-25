package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"kobili/user-ms/entities"
)

var (
	ErrTokenExpired = errors.New("the provided token has expired")
	ErrUserNotFound = errors.New("the user associated with the token could not be found")
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

/*
Retrieves the user given by the User ID in the subject claim of the given `token`.

# Returns

- `ErrTokenExpired` if the token has expired

- `ErrUserNotFound` if the user associated with the token does not exist
*/
func GetUserFromJWT(token string, dbConn *sql.DB, ctx context.Context) (*entities.UserEntity, error) {
	key := os.Getenv("JWT_SIGNING_KEY")
	if key == "" {
		return nil, fmt.Errorf("JWT_SIGNING_KEY not set")
	}

	parsedToken, err := jwt.Parse(
		token,
		func(t *jwt.Token) (interface{}, error) {
			return []byte(key), nil
		},
		jwt.WithValidMethods([]string{
			jwt.SigningMethodHS256.Alg(),
		}),
		jwt.WithIssuedAt(),
		jwt.WithIssuer("kobili/user-ms"),
	)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, fmt.Errorf("failed to verify JWT: %v", err)
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("failed to read JWT claims")
	}

	userId, err := claims.GetSubject()
	if err != nil {
		return nil, fmt.Errorf("failed to get userId from JWT claims: %v", err)
	}

	var user entities.UserEntity
	err = dbConn.QueryRowContext(
		ctx,
		`SELECT id, username, password FROM users WHERE id = $1`,
		userId,
	).Scan(&user.Id, &user.Username, &user.Password)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, ErrUserNotFound
		}
		return nil, fmt.Errorf("error fetching user from token: %v", err)
	}

	return &user, nil
}
