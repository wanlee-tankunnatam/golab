package router

import (
	"go-lab/handler"

	"github.com/gorilla/mux"
)

func RegisterRoutes() *mux.Router {
	router := mux.NewRouter()

	// Users
	router.HandleFunc("/post-user", handler.InsertUser).Methods("POST")

	// Stock
	router.HandleFunc("/stock/{id}", handler.GetStockById).Methods("GET") // path param /stock/1
	router.HandleFunc("/post-stock", handler.InsertStock).Methods("POST")
	router.HandleFunc("/put-stock", handler.UpdateStock).Methods("PUT") // query param ?id=1&qty=5

	return router
}
