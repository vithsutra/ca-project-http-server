package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/repository"
)

type adminHandler struct {
	adminRepo *repository.AdminRepo
}

func NewAdminHandler(adminRepo *repository.AdminRepo) *adminHandler {
	return &adminHandler{
		adminRepo,
	}
}

func (h *adminHandler) AdminLoginHandler(ctx echo.Context) error {
	adminLoginResponse, statusCode, err := h.adminRepo.AdminLogin(ctx)

	if err != nil {
		response := models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}

		ctx.JSON(int(statusCode), response)
		return err
	}

	response := models.SuccessResponse{
		Status:  "success",
		Message: "login successfull",
		Data:    adminLoginResponse,
	}

	ctx.JSON(int(statusCode), response)
	return nil
}
