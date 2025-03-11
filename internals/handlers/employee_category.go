package handlers

import (
	"errors"

	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/repository"
)

type employeeCategoryHandler struct {
	repo *repository.EmployeeCategoryRepo
}

func NewEmployeeCategoryHandler(repo *repository.EmployeeCategoryRepo) *employeeCategoryHandler {
	return &employeeCategoryHandler{
		repo,
	}
}

func (h *employeeCategoryHandler) CreateEmployeeCategoryHandler(ctx echo.Context) error {
	statusCode, err := h.repo.CreateEmployeeCategory(ctx)

	if err != nil {
		response := &models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}
		ctx.JSON(int(statusCode), response)
		return err
	}

	response := &models.SuccessResponse{
		Status:  "success",
		Message: "employee category created successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *employeeCategoryHandler) GetEmployeeCategoriesHandler(ctx echo.Context) error {
	employeeCategories, statusCode, err := h.repo.GetEmployeeCategories(ctx)
	if err != nil {
		response := &models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}
		ctx.JSON(int(statusCode), response)
		return err
	}

	if employeeCategories == nil {
		response := &models.ErrorResponse{
			Status: "error",
			Error:  "employee categories was empty",
		}
		ctx.JSON(404, response)
		return errors.New("employee categories not found")
	}

	response := &models.SuccessResponse{
		Status:  "success",
		Message: "employee categories fetched successfully",
		Data:    employeeCategories,
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *employeeCategoryHandler) DeleteEmployeeCategory(ctx echo.Context) error {
	statusCode, err := h.repo.DeleteEmployeeCategory(ctx)
	if err != nil {
		response := &models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}
		ctx.JSON(404, response)
		return err
	}

	response := &models.SuccessResponse{
		Status:  "success",
		Message: "employee category delete successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil
}
