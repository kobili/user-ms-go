package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"

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

	db := kdb.Connect()
	defer db.Close()

	http.HandleFunc("/", handlers.GreetingMessage)

	http.HandleFunc("/api/users/register", handlers.CreateUser(db))
	http.HandleFunc("/api/users/login", handlers.Login(db))

	fmt.Println("Listening on 127.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
