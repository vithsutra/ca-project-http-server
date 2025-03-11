package repository

import (
	"errors"
	"log"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/pkg/utils"
)

type RootRepo struct {
	dbRepo models.AdminInterface
}

func NewRootRepo(dbRepo models.AdminInterface) *RootRepo {
	return &RootRepo{
		dbRepo: dbRepo,
	}
}

func (repo *RootRepo) CreateAdmin(ctx echo.Context) (int32, error) {
	createAdminRequest := new(models.CreateAdminRequest)

	if err := ctx.Bind(createAdminRequest); err != nil {
		return 400, errors.New("invalid json body format")
	}

	validation := validator.New()

	if err := validation.RegisterValidation("password", utils.PasswordValidater); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.Struct(createAdminRequest); err != nil {
		return 400, errors.New("request body validation error")
	}

	hashedPassword, err := utils.HashPassword(createAdminRequest.Password)

	if err != nil {
		log.Println("Error occurred while generating the hashed password, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	admin := models.Admin{
		AdminId:     uuid.NewString(),
		Name:        createAdminRequest.Name,
		Dob:         createAdminRequest.Dob,
		Email:       createAdminRequest.Email,
		PhoneNumber: createAdminRequest.PhoneNumber,
		ProfileUrl:  "pending",
		Password:    hashedPassword,
		Position:    createAdminRequest.Position,
	}

	if err := repo.dbRepo.CreateAdmin(&admin); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	return 201, nil

}

func (repo *RootRepo) GetAllAdmins(ctx echo.Context) ([]*models.AdminResponse, int32, error) {

	adminResponses, err := repo.dbRepo.GetAllAdmins()

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error")
	}
	return adminResponses, 200, nil
}
