package handlers

import (
	"fmt"
	"net/http"
)

func GreetingMessage(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, fmt.Sprintf("Method %s not allowed", r.Method), http.StatusMethodNotAllowed)
		return
	}

	fmt.Fprint(w, "Hello!")
}
