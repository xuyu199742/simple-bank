package db

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
	"log"
	"os"
	"simple-bank/config"
	"testing"
)

var testQueries *Queries
var testDb *pgxpool.Pool

func TestMain(t *testing.M) {
	var err error
	cfg, err := config.LoadConfig("../../")
	if err != nil {
		log.Fatal("load config file failed", err)
	}
	testDb, err = pgxpool.New(context.Background(), cfg.DbSource)
	if err != nil {
		log.Fatal("can not connection postgres:", err)
	}
	testQueries = New(testDb)
	os.Exit(t.Run())
}
