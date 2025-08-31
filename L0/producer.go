package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/IBM/sarama"
	"github.com/gorilla/mux"
)

func RunProducer() {

	router := mux.NewRouter()
	router.HandleFunc("/order/{id}", getOrder).Methods("GET")
	router.HandleFunc("/order", createOrder).Methods("POST")

	log.Println("Server started at :3000")
	log.Fatal(http.ListenAndServe(":3000", router))
}

var cache = map[string]Order{}

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
	if router.Method != http.MethodPost {
		http.Error(writer, "invaild request method", http.StatusMethodNotAllowed)
		return
	}
	var order Order
	if err := json.NewDecoder(router.Body).Decode(&order); err != nil {
		http.Error(writer, "invalid json", http.StatusBadRequest)
		return
	}

	orderInBytes, err := json.Marshal(order)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
		return
	}

	err = PushOrderToQueue("orders", orderInBytes)
	if err != nil {
		log.Println(err)
		http.Error(writer, err.Error(), http.StatusInternalServerError)
	}

	response := map[string]interface{}{
		"success": true,
		"msg":     "Order placed successfully!",
	}

	writer.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(writer).Encode(response); err != nil {
		log.Println(err)
		http.Error(writer, "Error placing order", http.StatusInternalServerError)
		return
	}
}

func ConnectProducer(brokers []string) (sarama.SyncProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Return.Successes = true
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Retry.Max = 5

	return sarama.NewSyncProducer(brokers, config)
}

func PushOrderToQueue(topic string, message []byte) error {
	brokers := []string{"localhost:9092"}
	producer, err := ConnectProducer(brokers)
	if err != nil {
		return err
	}

	defer producer.Close()

	msg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.StringEncoder(message),
	}

	partition, offset, err := producer.SendMessage(msg)
	if err != nil {
		return err
	}

	log.Printf("Order is stored in topic(%s)/partition(%d)/offset(%d)\n",
		topic,
		partition,
		offset)

	return nil
}
