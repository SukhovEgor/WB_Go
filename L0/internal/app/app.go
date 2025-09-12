package app

import (
	"fmt"
	"log"
	"net/http"	
	"encoding/json"

	"test-task/internal/storage"

	"github.com/gorilla/mux"
)

type App struct {
	repository storage.Repository
}

func NewApp(connStr string) (*App, error) {
	app := &App{}
	err := app.repository.InitRepository(connStr)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return nil, err
	}
	return app, nil
}

func (a *App) HomeHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to the home page!")
}

func  (a *App)  GetOrderById(w http.ResponseWriter, r *http.Request) {
	order_uid := mux.Vars(r)["order_uid"]

	order, exist, err := a.repository.FindOrderById(order_uid)
	if !exist {
		fmt.Fprintf(w, "Order %v does not exist\n", order_uid)
		return
	} else if err != nil {
		log.Printf("Finding order by id is failed: %v", err)
		return
	}
	json_data, err := json.MarshalIndent(order, "", "\t")
	if err != nil {
		log.Printf("Failed to create json: %v", err)
	}
	fmt.Fprintf(w, "%s\n", json_data)
}

func (a *App) CreateOrders(w http.ResponseWriter, r *http.Request) {
	
}

/* func (a *App) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert")
	order := models.Order{
		OrderUID:          "test_order2",
		TrackNumber:       "test_trackNum",
		Entry:             "test_entry",
		Locale:            "ru",
		InternalSignature: "test_sign",
		CustomerID:        "test_customer",
		DeliveryService:   "test_deliverySer",
		Shardkey:          "test_SK",
		SmID:              1,
		DateCreated:       time.Now(),
		OofShard:          "1",
		Delivery: models.Delivery{
			OrderUID: "test_order2",
			Name:     "test_name",
			Phone:    "123456",
			Zip:      "test_zip",
			City:     "test_city",
			Address:  "test_address",
			Region:   "test_region",
			Email:    "test_email",
		},
		Payment: models.Payment{
			OrderUID:     "test_order2",
			Transaction:  "test_tx",
			RequestID:    "test_req",
			Currency:     "RUB",
			Provider:     "test_prov",
			Amount:       100,
			PaymentDt:    time.Now().Unix(),
			Bank:         "test_bank",
			DeliveryCost: 200,
			GoodsTotal:   300,
			CustomFee:    10,
		},
		Items: []models.Item{
			{
				ID:          1,
				OrderUID:    "test_order2",
				ChrtID:      11111,
				TrackNumber: "test_trackNum",
				Price:       100,
				Rid:         "test_Rid",
				Name:        "test_Name",
				Sale:        10,
				Size:        "test_size",
				TotalPrice:  1000,
				NmID:        1111111,
				Brand:       "test_Brand",
				Status:      0,
			},
		},
	}

	a.repository.InsertToDB(&order)
} */

/* func (a *App) Select(w http.ResponseWriter, r *http.Request) {
	// Query all authors
	rows, err := a.conn.Query(context.Background(), "SELECT * FROM orders")
	if err != nil {
		log.Fatalf("Error querying orders: %v", err)
	}
	defer rows.Close()

	var orders []models.Order
	for rows.Next() {
		var order models.Order
		if err := rows.Scan(order.OrderUID, order.TrackNumber, order.Entry,
			order.Locale, order.InternalSignature, order.CustomerID,
			order.DeliveryService, order.Shardkey, order.SmID,
			order.DateCreated, order.OofShard); err != nil {
			log.Fatalf("Error scanning row: %v", err)
		}
		orders = append(orders, order)
	}
	log.Println("Orders:", orders)
} */

func (a *App) Close() {
	a.repository.Close()
}
