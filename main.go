package main

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"simple-bank/api"
	db "simple-bank/db/sqlc"
)

const (
	dbSource = "postgres://root:root123@localhost:5432/simple_bank?sslmode=disable"
)

func main() {
	conn, err := pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("can not connection postgres:", err.Error())
	}

	server := api.NewServer(db.NewStore(conn))

	if err = server.Star(); err != nil {
		log.Fatal("can not start server", err.Error())
	}
}
