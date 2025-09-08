package storage

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Repository struct {
	pool *pgxpool.Pool
}

func (repository *Repository) InitRepository(connStr string) error {
	config, err := pgxpool.ParseConfig(connStr)
	if err != nil {
		log.Printf("Unable to parse config: %v", err)
		return err
	}

	repository.pool, err = pgxpool.NewWithConfig(context.Background(), config)
	if err != nil {
		log.Printf("Unable to connect to database: %v", err)
		return err
	}
	return nil
}

func (repository *Repository) InsertToDB(order *Order) {
	conn, err := repository.pool.Acquire(context.Background())
	if err != nil {
		log.Printf("Unable to get connection from the Pool: %v", err)
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Printf("Error starting transaction: %v", err)
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), insertOrder,
		order.OrderUID, order.TrackNumber, order.Entry,
		order.Locale, order.InternalSignature, order.CustomerID,
		order.DeliveryService, order.Shardkey, order.SmID,
		order.DateCreated, order.OofShard)
	if err != nil {
		log.Printf("Error inserting order: %v", err)
	}

	delivery := &order.Delivery
	_, err = tx.Exec(context.Background(), insertDelivery,
		order.OrderUID, delivery.Name, delivery.Phone,
		delivery.Zip, delivery.City, delivery.Address,
		delivery.Region, delivery.Email)
	if err != nil {
		log.Printf("Error inserting delivery: %v", err)
	}

	payment := &order.Payment
	_, err = tx.Exec(context.Background(), insertPayment,
		order.OrderUID, payment.Transaction, payment.RequestID,
		payment.Currency, payment.Provider, payment.Amount,
		payment.PaymentDt, payment.Bank, payment.DeliveryCost,
		payment.GoodsTotal, payment.CustomFee)
	if err != nil {
		log.Printf("Error inserting payment: %v", err)
	}

	for i := 0; i < len(order.Items); i++ {
		item := &order.Items[i]
		_, err = tx.Exec(context.Background(), insertItem,
			order.OrderUID, item.ChrtID, item.TrackNumber,
			item.Price, item.Rid, item.Name, item.Sale,
			item.Size, item.TotalPrice, item.NmID,
			item.Brand, item.Status,
		)
		if err != nil {
			log.Printf("Error inserting items: %v", err)
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}

	log.Println("Insert is completed")

}

func (repository *Repository) Close() {
	repository.pool.Close()
}
