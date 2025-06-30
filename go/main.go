package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	kdb "kobili/user-ms/db"
	"kobili/user-ms/handlers"
)

func initDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env found. Skipping...")
	}
}

func main() {
	initDotEnv()

	jwtSigningKey := os.Getenv("JWT_SIGNING_KEY")
	if jwtSigningKey == "" {
		log.Fatalf("JWT_SIGNING_KEY not set")
	}

	db := kdb.Connect()
	defer db.Close()

	http.HandleFunc("/", handlers.GreetingMessage)

	http.HandleFunc("/api/users/register", handlers.CreateUser(db))
	http.HandleFunc("/api/users/login", handlers.Login(db))
	http.HandleFunc("/api/users/me", handlers.GetUser(db))

	fmt.Println("Listening on 127.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
