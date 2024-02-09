package router

import (
	"github.com/gorilla/mux"

	"samurai/controller"
)

func Router() *mux.Router {
	router := mux.NewRouter()

	router.HandleFunc("/api/users", controller.InsertUser).Methods("POST")
	router.HandleFunc("/api/stations", controller.InsertStation).Methods("POST")
	router.HandleFunc("/api/trains", controller.InsertTrain).Methods("POST")
	router.HandleFunc("/api/stations", controller.ListAllStations).Methods("GET")
	// /api/stations/{station_id}/trains
	router.HandleFunc("/api/stations/{station_id}/trains", controller.ListTrainsAtStation).Methods("GET")
	// wallet
	router.HandleFunc("/api/wallets/{wallet_id}", controller.PrintWallet).Methods("GET")
	// /api/wallets/{wallet_id} put
	router.HandleFunc("/api/wallets/{wallet_id}", controller.InsertMoneyIntoWallet).Methods("PUT")

	router.HandleFunc("/api/tickets", controller.PurchaseTicket).Methods("POST")

	router.HandleFunc("/api/routes", controller.BestTicket).Methods("GET")

	// router.HandleFunc("/api/books", controller.InsertBook).Methods("POST")
	// router.HandleFunc("/api/books/{id}", controller.UpdateABook).Methods("PUT")
	// router.HandleFunc("/api/books/{id}", controller.GetABook).Methods("GET")
	// router.HandleFunc("/api/books", controller.SearchBooks).Methods("GET")
	// router.HandleFunc("/api/books/{id}", controller.DeleteABook).Methods("DELETE")


	return router
}