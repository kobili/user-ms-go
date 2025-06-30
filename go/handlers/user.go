package handlers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"golang.org/x/crypto/bcrypt"
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

		// TODO: insert into db
		w.Header().Add("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]string{
			"username": requestBody.Username,
			"password": string(passwordHash),
		})
	}
}
