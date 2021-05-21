package main

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/joho/godotenv"
	"io"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := pgx.Connect(context.Background(), os.Getenv("DB_URL"))
	if err != nil {
		log.Fatalln(err)
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		log.Fatalln(err)
	}
	defer tx.Commit(context.Background())
	lo := tx.LargeObjects()

	func() {
		_, err = lo.Create(context.Background(), 1)
		f, err := lo.Open(context.Background(), 1, pgx.LargeObjectModeWrite)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		lf, err := os.Open("/tmp/vstore-test/video1.mp4")
		if err != nil {
			log.Fatalln(err)
		}
		defer lf.Close()
		_, err = io.Copy(f, lf)
		if err != nil {
			log.Fatalln(err)
		}
	}()

	func() {
		_, err = lo.Create(context.Background(), 2)
		f, err := lo.Open(context.Background(), 2, pgx.LargeObjectModeWrite)
		if err != nil {
			log.Fatalln(err)
		}
		defer f.Close()
		lf, err := os.Open("/tmp/vstore-test/video2.mp4")
		if err != nil {
			log.Fatalln(err)
		}
		defer lf.Close()
		_, err = io.Copy(f, lf)
		if err != nil {
			log.Fatalln(err)
		}
	}()
}
