package main

import (
	"fmt"
	config "go-lab/config"
	"go-lab/router"
)

func main() {
	config.ConnectPostgres()
	router.RegisterRoutes()

	fmt.Println("Server running at test")

	// port := "80"
	// fmt.Printf("Server running at http://0.0.0.0:%s\n", port)

	// bind ทุก interface (0.0.0.0) เพื่อเรียกจาก IP ภายนอกได้
	// log.Fatal(http.ListenAndServe("128.199.97.209", nil))

	// fmt.Println("Server running at http://localhost:8080")
	// log.Fatal(http.ListenAndServe(":8080", nil))
}
