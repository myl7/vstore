package dao

import (
	"context"
	"github.com/jackc/pgx/v4"
	"log"
	"os"
)

func GetConn() *pgx.Conn {
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	return conn
}
