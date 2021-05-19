package dao

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

var conn *pgx.Conn

func GetConn() *pgx.Conn {
	if conn != nil {
		return conn
	}

	var err error
	conn, err = pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}

func CloseConn() error {
	if conn != nil {
		err := conn.Close(context.Background())
		if err != nil {
			return err
		}
	}
	return nil
}
