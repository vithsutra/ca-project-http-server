package repository

import (
	"errors"
	"log"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
)

type EmployeeCategoryRepo struct {
	dbRepo models.CategoryInterface
}

func NewEmployeeCategoryRepo(dbRepo models.CategoryInterface) *EmployeeCategoryRepo {
	return &EmployeeCategoryRepo{
		dbRepo: dbRepo,
	}
}

func (repo *EmployeeCategoryRepo) CreateEmployeeCategory(ctx echo.Context) (int32, error) {
	createEmployeeCategoryRequest := new(models.CreateEmployeeCategoryRequest)

	if err := ctx.Bind(createEmployeeCategoryRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(createEmployeeCategoryRequest); err != nil {
		return 400, errors.New("request body validation error")
	}

	categoryName := strings.ToLower(createEmployeeCategoryRequest.CategoryName)

	employeeCategoryExists, err := repo.dbRepo.CheckEmployeeCategoryExists(createEmployeeCategoryRequest.AdminId, categoryName)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if employeeCategoryExists {
		return 400, errors.New("employee category already exists")
	}

	employeeCategory := &models.EmployeeCategory{
		CategoryId:          uuid.NewString(),
		AdminId:             createEmployeeCategoryRequest.AdminId,
		CategoryName:        categoryName,
		CategoryDescription: createEmployeeCategoryRequest.CategoryDescription,
	}

	if err := repo.dbRepo.CreateEmployeeCategory(employeeCategory); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 201, nil
}

func (repo *EmployeeCategoryRepo) GetEmployeeCategories(ctx echo.Context) ([]*models.EmployeeCategoryResponse, int32, error) {
	adminId := ctx.Param("adminId")

	employeeCategories, err := repo.dbRepo.GetEmployeeCategories(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}
	return employeeCategories, 200, nil
}

func (repo *EmployeeCategoryRepo) DeleteEmployeeCategory(ctx echo.Context) (int32, error) {
	categoryId := ctx.Param("categoryId")

	categoryIdExists, err := repo.dbRepo.CheckEmployeeCategoryIdExists(categoryId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !categoryIdExists {
		return 400, errors.New("employee category id not exists")
	}

	if err := repo.dbRepo.DeleteEmployeeCategory(categoryId); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}
	return 200, nil
}
