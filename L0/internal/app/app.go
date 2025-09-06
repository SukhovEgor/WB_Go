package app

import (
	"context"
	"log"
	"net/http"
	"time"

	models "test-task/internal/repository"

	"github.com/jackc/pgx/v5"
)

type App struct {
	conn *pgx.Conn
}

func NewApp(connStr string) (*App, error) {

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v", err)
	}

	app := &App{
		conn: conn,
	}

	return app, nil
}

func (a *App) Insert(w http.ResponseWriter, r *http.Request) {
	log.Println("Insert")

	tx, err := a.conn.Begin(context.Background())
	if err != nil {
		log.Fatalf("Error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	// Check connection to Db and inserting
	order := models.Order{
		OrderUID:          "test_order",
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
	}

	delivery := models.Delivery{
		Name:    "test_name",
		Phone:   "123456",
		Zip:     "test_zip",
		City:    "test_city",
		Address: "test_address",
		Region:  "test_region",
		Email:   "test_email",
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO orders 
		(order_uid,	track_number, entry, 
		locale, internal_signature, customer_id,
		delivery_service, shardkey, sm_id, date_created, oof_shard)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, order.TrackNumber, order.Entry,
		order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID,
		order.DateCreated, order.OofShard)
	if err != nil {
		log.Fatalf("Error inserting order: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO deliveries 
		(order_uid, name, phone, 
		zip, city, address, 
		region, email)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`,
		order.OrderUID, delivery.Name, delivery.Phone,
		delivery.Zip, delivery.City, delivery.Address,
		delivery.Region, delivery.Email)
	if err != nil {
		log.Fatalf("Error inserting delivery: %v", err)
	}
	_, err = tx.Exec(context.Background(), `
		INSERT INTO payments
		(order_uid, transaction, request_id,
		currency, provider, amount, 
		payment_dt, bank, delivery_cost, 
		goods_total, custom_fee)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)`,
		order.OrderUID, "test_tx", "test_req", "RUB", "test_prov", 100, time.Now().Unix(), "test_bank", 200, 300, 10)
	if err != nil {
		log.Fatalf("Error inserting payment: %v", err)
	}

	_, err = tx.Exec(context.Background(), `
		INSERT INTO items (order_uid, chrt_id, track_number, price, rid, name, 
			sale, size, total_price, nm_id, brand, status)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12)`,
		order.OrderUID, 11111, "test_trackNum", 100, "test_Rid", "test_Name",
		10, "test_size", 1000, 1111111, "test_Brand", 0)
	if err != nil {
		log.Fatalf("Error inserting item: %v", err)
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Println("Insert is completed")

}

func (a *App) Select(w http.ResponseWriter, r *http.Request) {
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
}

func (a *App) Close() {
	a.conn.Close(context.Background())
}
