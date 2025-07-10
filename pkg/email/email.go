package email

import (
	"github.com/vithsutra/ca_project_http_server/internals/models"
	redisqueue "github.com/vithsutra/ca_project_http_server/pkg/reddis"
)

type EmailService struct {
	redisRepo *redisqueue.RedisQueueRepo
}

func NewEmailService(redisRepo *redisqueue.RedisQueueRepo) models.AdminEmailServiceInterface {
	return &EmailService{
		redisRepo: redisRepo,
	}
}

func (e *EmailService) SendEmail(data []byte) error {
	return e.redisRepo.SendEmail(data)
}
