package main

import (
	"log"
	"os"

	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/pkg/aws_s3"
	"github.com/vithsutra/ca_project_http_server/pkg/database"
	redisqueue "github.com/vithsutra/ca_project_http_server/pkg/reddis"

	"github.com/vithsutra/ca_project_http_server/repository"
)

func Start(dbConnPool *connection, awsS3Connection *s3Connection, redisConn *RedisQueueConnection) {
	e := echo.New()

	postgresRepo := database.NewPostgresRepo(dbConnPool.pool)
	awsS3Repo := aws_s3.NewAwsS3Repo(awsS3Connection.s3Client)
	redisRepo := redisqueue.NewRedisQueueRepo(redisConn.Client)

	rootRepo := repository.NewRootRepo(postgresRepo)
	adminRepo := repository.NewAdminRepo(postgresRepo, awsS3Repo, redisRepo)
	employeeCategoryRepo := repository.NewEmployeeCategoryRepo(postgresRepo)
	userRepo := repository.NewUserRepo(postgresRepo, awsS3Repo, redisRepo)

	InitHttpRoutes(e, rootRepo, adminRepo, employeeCategoryRepo, userRepo)

	if err := postgresRepo.Init(); err != nil {
		log.Fatalln("❌ Error initializing the database:", err)
	}

	log.Println("✅ Database initialized successfully")

	serverListenAddress := os.Getenv("SERVER_LISTEN_ADDRESS")
	if serverListenAddress == "" {
		log.Fatalln("❌ Please set SERVER_LISTEN_ADDRESS env variable")
	}

	if err := e.Start(serverListenAddress); err != nil {
		log.Fatalln("❌ Failed to start server:", err)
	}
}
