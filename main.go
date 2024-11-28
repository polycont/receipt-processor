package main

import (
	"net/http"

	api "github.com/polycont/receipt-processor/api/router/handlers"

	"github.com/gorilla/mux"
)

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/receipts/process", api.ProcessReceipts).Methods("POST")
	router.HandleFunc("/receipts/{id}/points", api.CalculatePoints).Methods("GET")

	http.ListenAndServe(":8080", router)
}
