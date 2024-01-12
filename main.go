package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"simple-bank/api"
	"simple-bank/config"
	db "simple-bank/db/sqlc"
)

func main() {

	cfg, err := config.LoadConfig(".")
	if err != nil {
		log.Fatal("load config file fail:", err.Error())
	}

	conn, err := pgxpool.New(context.Background(), cfg.DbSource)
	if err != nil {
		log.Fatal("can not connection postgres:", err.Error())
	}

	server := api.NewServer(db.NewStore(conn))

	if err = server.Star(cfg.ServerAddress); err != nil {
		log.Fatal("can not start server", err.Error())
	}
}
