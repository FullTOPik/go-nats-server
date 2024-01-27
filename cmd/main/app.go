package main

import (
	"context"
	_ "github.com/lib/pq"
	"go_nats-streaming_pg/internal/config"
	"go_nats-streaming_pg/internal/db"
	"go_nats-streaming_pg/internal/httpserver"
	"go_nats-streaming_pg/internal/streaming"
	"log"
	"os"
	"os/signal"
)

func main() {
	cfg := config.MustLoad()
	database, err := db.NewPG(context.Background(), cfg.Database)
	if err != nil {
		log.Fatalf("Error to create new pool, %+v", err)
		return
	}
	cache := db.CacheInit(database)
	server := httpserver.NewServer(cache)

	streaming.Publish(cfg.Stan)
	streaming.GetDataFromSteaming(database, cfg.Stan)

	done := make(chan bool, 1)
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, os.Interrupt)
	go func() {
		<-quit
		log.Println("Server is shutting down...")
		database.DeleteCacheState()
		server.FinishServer()
		done <- true
	}()
	<-done
}
