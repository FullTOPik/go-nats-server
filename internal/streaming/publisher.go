package streaming

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"go_nats-streaming_pg/internal/config"
	"go_nats-streaming_pg/internal/db"
	"log"
)

func Publish(nats config.Stan) {
	log.Println("Start publishing")
	sc, _ := stan.Connect(nats.ClusterID, nats.ClientID, stan.NatsURL(nats.SvrURL))
	defer sc.Close()
	item1 := db.Product{
		ChrtId:      1,
		TrackNumber: "TrackNumber",
		Price:       20,
		RId:         "RId",
		Name:        "Name",
		Sale:        10,
		Size:        "Size",
		TotalPrice:  20,
		NmId:        302010,
		Brand:       "Brand",
		Status:      1,
	}
	item2 := db.Product{
		ChrtId:      2,
		TrackNumber: "TrackNumber",
		Price:       31210,
		RId:         "RId",
		Name:        "Name",
		Sale:        112310,
		Size:        "L",
		TotalPrice:  20,
		NmId:        302010,
		Brand:       "guBrandsi",
		Status:      1,
	}
	delivery := db.Delivery{
		Name:    "Name",
		Phone:   "mobile",
		ZIP:     "3464",
		City:    "dMoskow",
		Address: "KK",
		Region:  "RR",
		Email:   "@",
	}
	payment := db.Payment{
		Transaction:  "ddsdf",
		RequestId:    "dasdfsds",
		Currency:     "dssdfa",
		Provider:     "ddfs",
		Amount:       10,
		PaymentDt:    11,
		Bank:         "alsdff",
		DeliveryCost: 12,
		GoodsTotal:   23,
		CustomFee:    32,
	}
	order := db.Order{
		OrderUID:          "Order 2",
		TrackNumber:       "track 2",
		Entry:             "adsdas",
		Delivery:          delivery,
		Payment:           payment,
		Items:             []db.Product{item1, item2},
		Locale:            "das",
		InternalSignature: "dddddd",
		CustomerId:        "customer 2",
		DeliveryService:   "DS2",
		Shardkey:          "adsasd",
		SmId:              123123,
		DateCreated:       "12.12.2013",
		CofShard:          "asdasd",
	}

	subject := nats.Subject
	for n := 0; n < 1; n++ {
		msg, _ := json.Marshal(order)
		sc.Publish(subject, msg)
		log.Println("Publish successfully")
	}
}
