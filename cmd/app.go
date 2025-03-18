package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/pkg/aws_s3"
	"github.com/vithsutra/ca_project_http_server/pkg/database"
	"github.com/vithsutra/ca_project_http_server/pkg/rabbitmq"
	"github.com/vithsutra/ca_project_http_server/repository"
)

func Start(dbConnPool *connection, awsS3Connection *s3Connection, rabbitmqConn *rabbitmqConnection) {
	e := echo.New()

	postgresRepo := database.NewPostgresRepo(dbConnPool.pool)

	awsS3Repo := aws_s3.NewAwsS3Repo(awsS3Connection.s3Client)

	rabbitmqRepo := rabbitmq.NewRabbitmqRepo(rabbitmqConn.conn, rabbitmqConn.chann)

	rootRepo := repository.NewRootRepo(postgresRepo)

	adminRepo := repository.NewAdminRepo(postgresRepo, awsS3Repo, rabbitmqRepo)

	employeeCategoryRepo := repository.NewEmployeeCategoryRepo(postgresRepo)

	userRepo := repository.NewUserRepo(postgresRepo, awsS3Repo, rabbitmqRepo)

	InitHttpRoutes(
		e,
		rootRepo,
		adminRepo,
		employeeCategoryRepo,
		userRepo,
	)

	if err := postgresRepo.Init(); err != nil {
		log.Fatalln("error occurred with database while initializing the database, Error: ", err.Error())
	}

	log.Println("database initialized successfully")

	serverListenAddres := os.Getenv("SERVER_LISTEN_ADDRESS")

	if serverListenAddres == "" {
		log.Fatalln("please set the SERVER_LISTEN_ADDRESS env variable")
	}
	if err := e.Start(serverListenAddres); err != nil {
		log.Fatalln("error occurred while starting the server: ", err)
	}
}
