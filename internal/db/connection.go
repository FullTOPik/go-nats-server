package db

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"go_nats-streaming_pg/internal/config"
	"log"
)

type Postgres struct {
	db        *pgxpool.Pool
	cacheInst *Cache
}

var pgInstance *Postgres

func (pg *Postgres) SetCache(cache *Cache) {
	pg.cacheInst = cache
}

func NewPG(ctx context.Context, dbconf config.Database) (*Postgres, error) {
	connString := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbconf.Host, dbconf.Port, dbconf.User, dbconf.Password, dbconf.DBName)
	db, err := pgxpool.New(ctx, connString)
	if err != nil {
		log.Fatalf("unable to create connection pool: %w", err)
	}

	pgInstance = &Postgres{db, &Cache{}}

	if err := pgInstance.Ping(context.Background()); err != nil {
		panic(err)
	}

	log.Println("Successfully connection to db")

	return pgInstance, nil
}

func (pg *Postgres) Ping(ctx context.Context) error {
	return pg.db.Ping(ctx)
}

func (pg *Postgres) Close() {
	pg.db.Close()
}
