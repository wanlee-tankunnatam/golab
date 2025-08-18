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

	// เอา router ที่ return จาก RegisterRoutes
	r := router.RegisterRoutes() // ✅ assign router

	fmt.Println("Server running at http ://128.199.97.209:80")
	log.Fatal(http.ListenAndServe(":80", r)) // ✅ ใช้ router แทน nil
}
