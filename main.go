package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"
	"os"

	kdb "kobili/user-ms/db"
)

func greetingMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	fmt.Fprint(w, "Hello!")
}

func initDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env found. Skipping...")
	}
}

func main() {
	initDotEnv()
	fmt.Println("POSTGRES_URL =", os.Getenv("POSTGRES_URL"))
	db := kdb.Connect()
	defer db.Close()

	http.HandleFunc("/", greetingMessage)

	fmt.Println("Listening on 127.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
