package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"net/http"

	kdb "kobili/user-ms/db"
)

func initDotEnv() {
	err := godotenv.Load()
	if err != nil {
		fmt.Println("No .env found. Skipping...")
	}
}

func greetingMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	fmt.Fprint(w, "Hello!")
}

func main() {
	initDotEnv()

	db := kdb.Connect()
	defer db.Close()

	http.HandleFunc("/", greetingMessage)

	fmt.Println("Listening on 127.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
