package main

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func init() {

	serverMode := os.Getenv("SERVER_MODE")

	if serverMode == "dev" {
		if err := godotenv.Load(); err != nil {
			log.Fatalln("missing the .env file", err)
			return
		}
		log.Println("running server in development mode..")
		log.Println(".env file loaded successfully..")
		return
	}

	if serverMode == "prod" {
		log.Println("running server in production mode..")
		return
	}

	log.Fatalln("please set SERVER_MODE to dev for development stage or to prod for production stage")

}

func main() {
	dbConnPool := NewDatabase()
	defer dbConnPool.CloseConnection()
	dbConnPool.CheckDatabase()

	awsS3Connection := NewS3Connection()

	rabbitmqConn := NewRabbitmqConnection()

	defer rabbitmqConn.conn.Close()
	defer rabbitmqConn.chann.Close()

	Start(dbConnPool, awsS3Connection, rabbitmqConn)
}
