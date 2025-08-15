package router

import (
	"net/http"

	"go-lab/handler"
)

func RegisterRoutes() {
	http.HandleFunc("/insert-user", handler.InsertUser)
	http.HandleFunc("/insert-stock", handler.InsertStock)
	http.HandleFunc("/update-stock", handler.UpdateStock)
}
