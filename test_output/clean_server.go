package main

import (
	"fmt"
	"log"
	"net/http"
)

func helloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	jsonResponse := `{"message": "Hello, World!"}`
	fmt.Fprintf(w, jsonResponse)
}

func main() {
	http.HandleFunc("/hello", helloHandler)

	fmt.Println("Server starting on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}