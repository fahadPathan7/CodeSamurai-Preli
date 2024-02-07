package router

import (
	"github.com/gorilla/mux"

	"samurai/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/books", controller.InsertBook).Methods("POST")
	router.HandleFunc("/api/books/{id}", controller.UpdateABook).Methods("PUT")
	router.HandleFunc("/api/books/{id}", controller.GetABook).Methods("GET")
	router.HandleFunc("/api/books", controller.SearchBooks).Methods("GET")


	return router
}