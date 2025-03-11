package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/repository"
)

type rootHandler struct {
	rootRepo *repository.RootRepo
}

func NewRootHandler(rootRepo *repository.RootRepo) *rootHandler {
	return &rootHandler{
		rootRepo,
	}
}

func (h *rootHandler) CreateAdminHandler(ctx echo.Context) error {

	statusCode, err := h.rootRepo.CreateAdmin(ctx)

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
		Message: "admin created successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *rootHandler) GetAllAdminsHandler(ctx echo.Context) error {
	adminResponses, statusCode, err := h.rootRepo.GetAllAdmins(ctx)
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
		Message: "admins fetched successfully",
		Data:    adminResponses,
	}
	ctx.JSON(int(statusCode), response)
	return nil
}
