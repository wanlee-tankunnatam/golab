package main

import (
	"fmt"
	"log"
	"net/http"

	config "go-lab/config"
	"go-lab/router"
)

func main() {
	config.ConnectPostgres()
	router.RegisterRoutes()

	fmt.Println("Server running at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
