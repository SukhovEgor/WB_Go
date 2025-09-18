package app

import (
	"encoding/json"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	cache "test-task/internal/cache"
	"test-task/internal/kafka"
	models "test-task/internal/models"
	"test-task/internal/storage"

	"github.com/IBM/sarama"
	"github.com/brianvoe/gofakeit/v7"
	"github.com/gorilla/mux"
)

type App struct {
	repository storage.Repository
	cache      cache.Cache
	Producer   sarama.SyncProducer
	Consumer   sarama.Consumer
}

func NewApp(connStr string) (*App, error) {
	app := &App{}

	err := app.repository.InitRepository(connStr)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return nil, err
	}

	brokers := []string{"localhost:9092"}

	producer, err := kafka.ConnectProducer(brokers)
	if err != nil {
		log.Printf("Unable to init Kafka producer: %v", err)
		return nil, err
	}
	app.Producer = producer

	consumer, err := kafka.ConnectConsumer(brokers)
	if err != nil {
		log.Printf("Unable to init Kafka consumer: %v", err)
		return nil, err
	}
	app.Consumer = consumer

	return app, nil
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	html, err := os.ReadFile("../web/index.html")
	if err != nil {
		log.Printf("Error reading index.html: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "internal error")
		return
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(html)
}

func (a *App) GetOrderById(w http.ResponseWriter, r *http.Request) {
	order_uid := mux.Vars(r)["order_uid"]
	log.Printf("Searching : %v", order_uid)

	order, exist, err := a.repository.FindOrderById(order_uid)

	if !exist {
		fmt.Fprintf(w, "Order %v does not exist\n", order_uid)
		return
	} else if err != nil {
		log.Printf("Finding order by id is failed: %v", err)
		return
	}

	kafka.DoRequest(a.Producer, a.Consumer, order_uid, "get_order_by_id", "get_order_by_id_response")

	json_data, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		log.Printf("Failed to create json: %v", err)
	}
	fmt.Fprintf(w, "%s\n", json_data)
}

func (a *App) HandleGetOrderByID(uid string) (interface{}, error) {
	uid = strings.Trim(uid, `"`)
	log.Printf("HandleSearching : %v", uid)
	order, exist, err := a.repository.FindOrderById(uid)
	if err != nil {
		log.Printf("DB fetch error: %v", err)
		return nil, err
	}
	if !exist {
		return nil, fmt.Errorf("Order %s is not found", uid)
	}
	return order, nil
}

func (a *App) HandleCreateOrders(data string) (interface{}, error) {
	orderCount, err := strconv.Atoi(data)
	if err != nil {
		log.Printf("Parse error: %v", err)
		return nil, err
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	ordersAdded := 0

	var orders []models.Order
	for i := 0; i < orderCount; i++ {
		order, err := createRandomOrder(rng)

		if err != nil {
			return nil, err
		}

		if err := a.repository.InsertToDB(&order); err != nil {
			log.Printf("DB inserting error: %v", err)
			return nil, err
		}
		ordersAdded++
		orders = append(orders, order)
	}

	return orders, nil
}

func (a *App) CreateOrders(w http.ResponseWriter, r *http.Request) {
	orderCount := 2
	msg := kafka.DoRequest(a.Producer, a.Consumer, orderCount,
		"post_order", "post_order_response")

	var orders []models.Order
	if err := json.Unmarshal([]byte(msg), &orders); err != nil {
		response := map[string]interface{}{
			"error":   true,
			"message": msg,
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(response)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		log.Printf("Error while creating response: %v", err)
	}
}

func createRandomOrder(rng *rand.Rand) (models.Order, error) {
	var order models.Order
	var delivery models.Delivery
	var payment models.Payment

	if err := gofakeit.Struct(&order); err != nil {
		return order, err
	}
	if err := gofakeit.Struct(&delivery); err != nil {
		return order, err
	}
	if err := gofakeit.Struct(&payment); err != nil {
		return order, err
	}

	itemCount := rng.Intn(10) + 1
	items := make([]models.Item, 0, itemCount)
	var goodsTotal float64

	for i := 0; i < itemCount; i++ {
		var item models.Item
		if err := gofakeit.Struct(&item); err != nil {
			return order, err
		}

		quantity := rng.Intn(5) + 1

		item.Sale = rng.Intn(51)

		item.TotalPrice = float64(item.Price*quantity) * (1 - float64(item.Sale)/100.0)

		goodsTotal += item.TotalPrice
		items = append(items, item)
	}

	payment.GoodsTotal = goodsTotal
	payment.Amount = payment.DeliveryCost + goodsTotal + payment.CustomFee

	order.Delivery = delivery
	order.Payment = payment
	order.Items = items

	return order, nil
}

func (a *App) Close() {
	a.repository.Close()
	if a.Producer != nil {
		a.Producer.Close()
	}
	if a.Consumer != nil {
		a.Consumer.Close()
	}
}
