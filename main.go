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

	port := "80"
	fmt.Printf("Server running at http://0.0.0.0:%s\n", port)

	// bind ทุก interface (0.0.0.0) เพื่อเรียกจาก IP ภายนอกได้
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))

	// fmt.Println("Server running at http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
