package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"kobili/user-ms/entities"
	"kobili/user-ms/utils"
)

const BCRYPT_COST = 3

type CreateUserRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func CreateUser(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing request body: %v", err), http.StatusInternalServerError)
			return
		}

		var requestBody CreateUserRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error parsing JSON: %v", err), http.StatusBadRequest)
			return
		}

		passwordHash, err := bcrypt.GenerateFromPassword([]byte(requestBody.Password), BCRYPT_COST)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error hashing password: %v", err), http.StatusInternalServerError)
			return
		}

		var user entities.UserEntity
		err = dbConn.QueryRowContext(
			r.Context(),
			`INSERT INTO "users" ("username", "password")
			VALUES ($1, $2)
			RETURNING "id", "username", "password"`,
			requestBody.Username,
			string(passwordHash),
		).Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			http.Error(w, fmt.Sprintf("Error writing to db: %v", err), http.StatusInternalServerError)
			return
		}

		token, err := utils.CreateJWTForUser(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to sign JWT: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"id":       user.Id,
			"username": user.Username,
			"token":    token,
		})
	}
}

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func Login(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}

		body, err := io.ReadAll(r.Body)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to read request body: %v", err), http.StatusInternalServerError)
			return
		}

		var requestBody LoginRequest
		err = json.Unmarshal(body, &requestBody)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to parse JSON: %v", err), http.StatusBadRequest)
			return
		}

		var user entities.UserEntity
		err = dbConn.QueryRowContext(
			r.Context(),
			`SELECT id, username, password FROM users WHERE username = $1`,
			requestBody.Username,
		).Scan(&user.Id, &user.Username, &user.Password)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Invalid username or password", http.StatusNotFound)
				return
			}
			http.Error(w, fmt.Sprintf("Failed to retreive user: %v", err), http.StatusInternalServerError)
			return
		}

		err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(requestBody.Password))
		if err != nil {
			http.Error(w, "Invalid username or password", http.StatusNotFound)
			return
		}

		signedToken, err := utils.CreateJWTForUser(user)
		if err != nil {
			http.Error(w, fmt.Sprintf("Failed to sign JWT: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"token": signedToken,
		})
	}
}

func GetUser(dbConn *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
			return
		}

		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			http.Error(w, "Missing Authorization header", http.StatusUnauthorized)
			return
		}

		authHeaderPieces := strings.Split(authHeader, " ")
		if len(authHeaderPieces) != 2 || authHeaderPieces[0] != "Bearer" {
			http.Error(w, "Authorization header should be a Bearer token", http.StatusUnauthorized)
			return
		}

		accessToken := authHeaderPieces[1]
		// jwt.WithValidMethods([]string{
		// 	jwt.SigningMethodES256.Name,
		// })

		user, err := utils.GetUserFromJWT(accessToken, dbConn, r.Context())
		if err != nil {
			// TODO: make this error handling less general
			http.Error(w, fmt.Sprintf("Error retrieving user from JWT: %v", err), http.StatusInternalServerError)
			return
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]string{
			"id":       user.Id,
			"username": user.Username,
		})
	}
}
