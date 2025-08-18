package main

import (
	"fmt"
	config "go-lab/config"
	"go-lab/router"
	"log"
	"net/http"
)

func main() {
	config.ConnectPostgres()

	// เอา router ที่ return มาจาก RegisterRoutes
	r := router.RegisterRoutes()

	fmt.Println("Server running at IP http://128.199.97.209:80")
	log.Fatal(http.ListenAndServe(":80", r)) // ใส่ router แทน nil
}
