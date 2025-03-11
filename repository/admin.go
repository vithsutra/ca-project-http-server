package repository

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/pkg/utils"
)

type AdminRepo struct {
	dbRepo models.AdminInterface
}

func NewAdminRepo(dbRepo models.AdminInterface) *AdminRepo {
	return &AdminRepo{
		dbRepo: dbRepo,
	}
}

func (repo *AdminRepo) AdminLogin(ctx echo.Context) (*models.AdminLoginResponse, int32, error) {
	adminLoginRequest := new(models.AdminLoginRequest)

	if err := ctx.Bind(adminLoginRequest); err != nil {
		return nil, 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(adminLoginRequest); err != nil {
		return nil, 400, errors.New("request body validation error")
	}

	adminEmailsExists, err := repo.dbRepo.CheckAdminEmailsExists(adminLoginRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if !adminEmailsExists {
		return nil, 400, errors.New("admin email not exists")
	}

	adminId, userName, hashedPassword, err := repo.dbRepo.GetAdminForLogin(adminLoginRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if err := utils.CheckPassword(hashedPassword, adminLoginRequest.Password); err != nil {
		return nil, 401, errors.New("incorrect password")
	}

	token, err := utils.GenerateToken(adminId, adminLoginRequest.Email, userName)

	if err != nil {
		log.Println("error occurred while generating the token, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	return &models.AdminLoginResponse{
		Token: token,
	}, 200, nil

}
