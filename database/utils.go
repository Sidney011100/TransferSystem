package database

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
)

const (
	ErrDupKey = "ERROR: duplicate key value violates unique constraint"
)

var conn *pgx.Conn

func InitDatabase(dsn string) {
	var err error
	conn, err = pgx.Connect(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
}

func CloseDatabase() {
	if conn != nil {
		conn.Close(context.Background())
	}
}
