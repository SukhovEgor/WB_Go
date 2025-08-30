package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Order struct {
	OrderUID   string `json:"order_uid"`
	CustomerID string `json:"customer_id"`
}

var cache = map[string]Order{}

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", getOrder).Methods("GET")
	router.HandleFunc("/order", createOrder).Methods("POST")

	log.Println("Server started at :8081")
	log.Fatal(http.ListenAndServe(":8081", router))
}

func getOrder(writer http.ResponseWriter, router *http.Request) {
	vars := mux.Vars(router)
	id := vars["id"]

	if order, ok := cache[id]; ok {
		json.NewEncoder(writer).Encode(order)
		return
	}

	http.Error(writer, "order not found", http.StatusNotFound)
}

func createOrder(writer http.ResponseWriter, router *http.Request) {
	var order Order
	if err := json.NewDecoder(router.Body).Decode(&order); err != nil {
		http.Error(writer, "invalid json", http.StatusBadRequest)
		return
	}

	cache[order.OrderUID] = order
	writer.WriteHeader(http.StatusCreated)
	json.NewEncoder(writer).Encode(order)
}
