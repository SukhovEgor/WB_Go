package main

import (
	//"context"
	"test-task/internal/app"
	"log"
	"net/http"

	//"github.com/jackc/pgx/v5/pgxpool"
	"github.com/gorilla/mux"
)

func main() {

	//go RunProducer()
	//RunConsumer()

	// Подключение к БД
	connStr := "postgres://postgres:qwerty@localhost:5433/WB_ordersDB"
	newApp, err := app.NewApp(connStr)
	if err != nil {
		log.Fatalf("Failed to initialize")
	}
	defer newApp.Close()
	
	r := mux.NewRouter()

	r.HandleFunc("/api/add", newApp.Insert)
	//r.HandleFunc("/api/get", newApp.Select)
	log.Println("Server started at :3000")
	http.ListenAndServe(":3000", r)

}
