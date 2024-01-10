package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"testing"
)

const (
	dbSource = "postgres://root:root123@localhost:5432/simple_bank?sslmode=disable"
)

var testQueries *Queries
var testDb *pgxpool.Pool

func TestMain(t *testing.M) {
	var err error
	testDb, err = pgxpool.New(context.Background(), dbSource)
	if err != nil {
		log.Fatal("can not connection postgres:", err)
	}
	testQueries = New(testDb)
	os.Exit(t.Run())
}
