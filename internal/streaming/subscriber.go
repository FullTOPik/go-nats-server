package streaming

import (
	"encoding/json"
	"github.com/nats-io/stan.go"
	"go_nats-streaming_pg/internal/config"
	"go_nats-streaming_pg/internal/db"
	"log"
)

func GetDataFromSteaming(database *db.Postgres, stanconf config.Stan) {
	log.Println("Subscribing")
	sc, _ := stan.Connect(stanconf.ClusterID, stanconf.ClientID, stan.NatsURL(stanconf.SvrURL),
		stan.SetConnectionLostHandler(func(_ stan.Conn, reason error) {
			log.Fatalf("connection lost, reason %v", reason)
		}))

	subject := stanconf.Subject
	var newOrder db.Order
	_, err := sc.Subscribe(subject, func(msg *stan.Msg) {
		log.Println("Get new message")
		err := json.Unmarshal(msg.Data, &newOrder)
		if err != nil {
			log.Println("Bad format message")
			msg.Ack()
		}

		err = database.SaveOrder(&newOrder)
		if err != nil {
			log.Println("Error to save order")
		} else {
			msg.Ack()
		}
	}, stan.DeliverAllAvailable(), stan.DurableName("test"), stan.SetManualAckMode())

	if err != nil {
		log.Fatal(err)
	}
}
