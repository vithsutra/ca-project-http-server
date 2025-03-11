package main

import (
	"context"
	"log"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

type connection struct {
	pool *pgxpool.Pool
}

func NewDatabase() *connection {
	dbUrl := os.Getenv("DB_URL")

	if dbUrl == "" {
		log.Fatalln("DB_URL env variable is missing")
	}
	pool, err := pgxpool.New(context.Background(), dbUrl)
	if err != nil {
		log.Fatalln("error occurred while connecting to database, Error: ", err.Error())
	}

	log.Println("connected to the database with pooling")

	return &connection{
		pool: pool,
	}

}

func (db *connection) CheckDatabase() {
	if err := db.pool.Ping(context.Background()); err != nil {
		log.Fatalln("error occurred while performing database healthcheck, Error: ", err.Error())
	}

	log.Println("database was working correctly")
}

func (db *connection) CloseConnection() {
	db.pool.Close()
}
