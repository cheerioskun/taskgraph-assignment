package main

import (
	"fmt"
	"log"
	"net/http"
)

func main() {
	mux := setupRoutes()

	port := "8080"
	fmt.Printf("Server starting on port %s\n", port)

	if err := http.ListenAndServe("0.0.0.0:"+port, mux); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}
