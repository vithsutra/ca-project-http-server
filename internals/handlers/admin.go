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

func (h *adminHandler) GetAdminProfileDetailsHandler(ctx echo.Context) error {
	details, statusCode, err := h.adminRepo.GetAdminProfileDetails(ctx)

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
		Message: "admin profile details fetched successfully",
		Data:    details,
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *adminHandler) UpdateAdminNewPasswordHandler(ctx echo.Context) error {
	statusCode, err := h.adminRepo.UpdateAdminNewPassword(ctx)
	if err != nil {
		response := models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}

		ctx.JSON(int(statusCode), response)
		return err
	}

	response := models.SuccessResponse{
		Status:  "successs",
		Message: "admin new password updated successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil

}

func (h *adminHandler) UpdateAdminProfileInfoHandler(ctx echo.Context) error {
	statusCode, err := h.adminRepo.UpdateAdminProfileInfo(ctx)

	if err != nil {
		response := models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}

		ctx.JSON(int(statusCode), response)
		return err
	}

	response := models.SuccessResponse{
		Status:  "successs",
		Message: "admin profile updated successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *adminHandler) UpdateAdminProfilePictureHandler(ctx echo.Context) error {
	statusCode, err := h.adminRepo.UpdateProfilePicture(ctx)
	if err != nil {
		response := models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}

		ctx.JSON(int(statusCode), response)
		return err
	}

	response := models.SuccessResponse{
		Status:  "successs",
		Message: "admin profile picture updated successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *adminHandler) DeleteAdminProfilePictureHandler(ctx echo.Context) error {
	statusCode, err := h.adminRepo.DeleteProfilePicture(ctx)

	if err != nil {
		response := models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}

		ctx.JSON(int(statusCode), response)
		return err
	}

	response := models.SuccessResponse{
		Status:  "successs",
		Message: "admin profile picture deleted successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil

}
