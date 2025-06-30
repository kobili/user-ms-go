package main

import (
	"fmt"
	"net/http"
)

func greetingMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		fmt.Fprintf(w, "Method %s not allowed", r.Method)
		return
	}

	fmt.Fprint(w, "Hello!")
}

func main() {
	http.HandleFunc("/", greetingMessage)

	fmt.Println("Listening on 127.0.0.0:8080")
	http.ListenAndServe(":8080", nil)
}
