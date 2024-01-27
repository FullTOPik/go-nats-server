package db

import (
	"context"
	"github.com/jackc/pgx/v5"
	"log"

	_ "github.com/jackc/pgx/v5/pgxpool"
)

func (pg *Postgres) SaveOrder(order *Order) error {
	var (
		lastInsertId int64
		orderId      int64
		deliveryId   int64
		paymentId    int64
		itemsId      []int64
	)

	tx, err := pg.db.Begin(context.Background())
	if err != nil {
		log.Println("Error to create connection pool")
		return err
	}

	defer tx.Rollback(context.Background())

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO deliveries ("name", "phone", "zip", "city", "address", "region", "email") 
		VALUES ($1, $2, $3, $4, $5, $6, $7) RETURNING id`,
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.ZIP,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error to add in delivery")
		return err
	}

	deliveryId = lastInsertId

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO payments (transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) RETURNING id`,
		&order.Payment.Transaction,
		&order.Payment.RequestId,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	).Scan(&lastInsertId)
	if err != nil {
		log.Println("Error to add in payment")
		return err
	}

	paymentId = lastInsertId

	for _, item := range order.Items {
		err = tx.QueryRow(
			context.Background(),
			`INSERT INTO products (chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status)
				VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11) RETURNING id`,
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.RId,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		).Scan(&lastInsertId)
		if err != nil {
			log.Println("Error to add in item")
			return err
		}

		itemsId = append(itemsId, lastInsertId)
	}

	err = tx.QueryRow(
		context.Background(),
		`INSERT INTO orders (order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, shardkey, sm_id, cof_shard, date_created)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13) RETURNING id`,
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		deliveryId,
		paymentId,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.CofShard,
		&order.DateCreated,
	).Scan(&lastInsertId)

	if err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Error to add order, err = %+v", err)
			return err
		}
	}

	orderId = lastInsertId

	for _, itemId := range itemsId {
		_, err = tx.Exec(context.Background(), `INSERT INTO orderProducts (order_id, item_id)
		VALUES ($1, $2)`, orderId, itemId)
		if err != nil {
			log.Printf("Error to add in order_item %+v", err)
			return err
		}
	}

	log.Println("Order added in DB")

	pg.cacheInst.SetOrderInCache(*order, orderId)

	err = tx.Commit(context.Background())
	if err != nil {
		log.Printf("commit error, ERROR %+v", err)
		return err
	}

	return nil
}

func (pg *Postgres) GetOrderById(orderId int64) (Order, error) {
	var (
		order      Order
		deliveryId int64
		paymentId  int64
		itemId     int64
		item       Product
	)

	tx, err := pg.db.Begin(context.Background())
	if err != nil {
		log.Println("Error to connect db")
		return order, err
	}

	defer tx.Rollback(context.Background())
	err = tx.QueryRow(context.Background(), `SELECT order_uid, track_number, entry, delivery_id, payment_id, locale, internal_signature, customer_id, delivery_service, 
		shardkey, sm_id, cof_shard, date_created FROM orders WHERE id = $1`, orderId).Scan(
		&order.OrderUID,
		&order.TrackNumber,
		&order.Entry,
		&deliveryId,
		&paymentId,
		&order.Locale,
		&order.InternalSignature,
		&order.CustomerId,
		&order.DeliveryService,
		&order.Shardkey,
		&order.SmId,
		&order.CofShard,
		&order.DateCreated,
	)

	if err != nil {
		log.Printf("Failure to select order %+v", err)
		return order, err
	}

	err = tx.QueryRow(context.Background(), `SELECT name, phone, zip, city, address, region, email FROM deliveries WHERE id = $1`, deliveryId).Scan(
		&order.Delivery.Name,
		&order.Delivery.Phone,
		&order.Delivery.ZIP,
		&order.Delivery.City,
		&order.Delivery.Address,
		&order.Delivery.Region,
		&order.Delivery.Email,
	)

	if err != nil {
		log.Printf("Error to select delivery %+v", err)
		return order, err
	}

	err = tx.QueryRow(context.Background(), `SELECT transaction, request_id, currency, provider, amount, payment_dt, bank, delivery_cost, goods_total, custom_fee FROM payments WHERE id = $1`, paymentId).Scan(
		&order.Payment.Transaction,
		&order.Payment.RequestId,
		&order.Payment.Currency,
		&order.Payment.Provider,
		&order.Payment.Amount,
		&order.Payment.PaymentDt,
		&order.Payment.Bank,
		&order.Payment.DeliveryCost,
		&order.Payment.GoodsTotal,
		&order.Payment.CustomFee,
	)

	if err != nil {
		log.Printf("Error to select payments %+v", err)
		return order, err
	}

	allItems, err := tx.Query(context.Background(), `SELECT item_id FROM orderProducts WHERE order_id = $1`, orderId)

	if err != nil {
		log.Printf("Error to select all items %+v", err)
		return order, err
	}

	defer allItems.Close()

	conn, err := pg.db.Begin(context.Background())
	defer conn.Rollback(context.Background())

	if err != nil {
		log.Println("Error to get items on db")
	}

	for allItems.Next() {
		if err = allItems.Scan(&itemId); err != nil {
			log.Println("Error to search items")
			return order, err
		}
		err = conn.QueryRow(context.Background(), `SELECT chrt_id, track_number, price, rid, name, sale, size, total_price, nm_id, brand, status FROM products WHERE id = $1`, itemId).Scan(
			&item.ChrtId,
			&item.TrackNumber,
			&item.Price,
			&item.RId,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmId,
			&item.Brand,
			&item.Status,
		)

		if err != nil {
			log.Printf("Error to select items %+v", err)
			return order, err
		}

		order.Items = append(order.Items, item)
	}

	return order, nil
}

func (pg *Postgres) GetCacheStateInDB() (map[int64]Order, error) {
	item := make(map[int64]Order)
	tx, err := pg.db.Begin(context.Background())
	if err != nil {
		log.Println("Error to connect db")
		return item, err
	}
	defer tx.Rollback(context.Background())

	cacheState, err := tx.Query(context.Background(), `SELECT order_id FROM cache_state`)
	if err != nil {
		log.Println("Error to select cache from db")
		return item, err
	}

	defer cacheState.Close()
	var orderId int64

	for cacheState.Next() {
		if err := cacheState.Scan(&orderId); err != nil {
			log.Println("Unable to get order_id")
			return item, err
		}

		order, err := pg.GetOrderById(orderId)
		if err != nil {
			log.Println("Error to get order in db")
			return item, err
		}

		item[orderId] = order
	}

	if len(item) == 0 {
		log.Println("Cache empty")
		return item, nil
	}

	log.Println("Cache state restored")
	return item, nil
}

func (pg *Postgres) SaveCacheInDB(orderId int64) {
	err := pg.db.QueryRow(context.Background(), `INSERT INTO cache_state (order_id) VALUES ($1)`, orderId).Scan()

	if err != nil {
		if err != pgx.ErrNoRows {
			log.Printf("Error in save cache, %+v", err)
		}
	}

	log.Println("Cache state successfully save in db")
}

func (pg *Postgres) DeleteCacheState() {
	log.Println("Start delete cache state")
	_, err := pg.db.Exec(context.Background(), "DELETE FROM cache_state")
	if err != nil {
		log.Printf("Error to delete cache state, %+v", err)
		return
	}
	log.Println("Successfully delete cache state")
}
