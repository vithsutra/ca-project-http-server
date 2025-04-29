package handlers

import (
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/repository"
)

type userHandler struct {
	repo *repository.UserRepo
}

func NewUserHandler(repo *repository.UserRepo) *userHandler {
	return &userHandler{
		repo,
	}
}

func (h *userHandler) CreateUserHandler(ctx echo.Context) error {
	userId, statusCode, err := h.repo.CreateUser(ctx)

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
		Message: "user created successfully",
		Data: map[string]string{
			"user_id": userId},
	}

	ctx.JSON(int(statusCode), response)
	return nil
}
func (h *userHandler) GetUserProfileDetailsHandler(ctx echo.Context) error {
	details, statusCode, err := h.repo.GetUserProfileDetails(ctx)

	if err != nil {
		response := &models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}
		ctx.JSON(int(statusCode), response)
		return err
	}

	response := models.SuccessResponse{
		Status:  "success",
		Message: "user profile details fetched successfull",
		Data:    details,
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *userHandler) GetUsers(ctx echo.Context) error {
	usersReponse, statusCode, err := h.repo.GetUsers(ctx)

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
		Message: "users fetched successfully",
		Data:    usersReponse,
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) DeleteUser(ctx echo.Context) error {
	statusCode, err := h.repo.DeleteUser(ctx)

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
		Message: "user deleted successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) UserLoginHandler(ctx echo.Context) error {
	userLoginResponse, statusCode, err := h.repo.UserLogin(ctx)

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
		Message: "user login successfull",
		Data:    userLoginResponse,
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *userHandler) UserWorkLoginHandler(ctx echo.Context) error {
	statusCode, err := h.repo.UserWorkLogin(ctx)

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
		Message: "user work login successfull",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) UserWorkLogoutHandler(ctx echo.Context) error {
	statusCode, err := h.repo.UserWorkLogout(ctx)
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
		Message: "user work logout successfull",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) ApplyUserLeaveHandler(ctx echo.Context) error {
	statusCode, err := h.repo.ApplyUserLeave(ctx)
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
		Message: "user leave apllied successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) GetUserPendingLeavesHandler(ctx echo.Context) error {
	pendingLeavesCount, pendingLeaves, statusCode, err := h.repo.GetAllUsersPendingLeaves(ctx)

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
		Message: "users pending leaves fetched successfully",
		Data: map[string]interface{}{
			"total_count": pendingLeavesCount,
			"leaves":      pendingLeaves,
		},
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *userHandler) GetUserLeavesHandler(ctx echo.Context) error {
	usersLeavesCount, userLeaves, statusCode, err := h.repo.GetUserLeaves(ctx)

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
		Message: "user leaves fetched successfully",
		Data: map[string]interface{}{
			"total_count": usersLeavesCount,
			"laves":       userLeaves,
		},
	}

	ctx.JSON(int(statusCode), response)
	return nil

}

func (h *userHandler) CancelUserLeaveHandler(ctx echo.Context) error {
	statusCode, err := h.repo.CancelUserLeave(ctx)

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
		Message: "user leave canceled successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) GrantUserLeaveHandler(ctx echo.Context) error {
	statusCode, err := h.repo.GrantUserLeave(ctx)

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
		Message: "user leave granted successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *userHandler) UserProfileInfoUpdateHandler(ctx echo.Context) error {
	statusCode, err := h.repo.UpdateUserProfileInfo(ctx)

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
		Message: "user profile info updated successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil

}

func (h *userHandler) UpdateUserProfilePictureHandler(ctx echo.Context) error {
	statusCode, err := h.repo.UpdateUserProfilePicture(ctx)

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
		Message: "user profile picture updated successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil

}

func (h *userHandler) DeleteProfilePictureHandler(ctx echo.Context) error {
	statusCode, err := h.repo.DeleteUserProfilePicture(ctx)

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
		Message: "user profile picture deleted successfully",
	}

	ctx.JSON(int(statusCode), response)

	return nil

}

func (h *userHandler) GetUserLastProfileUpdateTimeHandler(ctx echo.Context) error {
	lastUpdateTime, statusCode, err := h.repo.GetUserLastProfileUpdateTime(ctx)

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
		Message: "user last profile update time fetched successfully",
		Data: &models.UserLastProfileUpdateTimeResponse{
			Time: lastUpdateTime,
		},
	}

	ctx.JSON(int(statusCode), response)

	return nil
}

func (h *userHandler) UpdateUserNewPaswordHandler(ctx echo.Context) error {
	statusCode, err := h.repo.UpdateUserNewPassword(ctx)

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
		Message: "user password updated successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) UserForgotPasswordHandler(ctx echo.Context) error {
	statusCode, err := h.repo.UserForgotPassword(ctx)
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
		Message: "otp sent successfully",
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) ValidateUserOtpHandler(ctx echo.Context) error {
	token, statusCode, err := h.repo.ValidateUserOtp(ctx)
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
		Message: "otp validation successfull",
		Data: &models.UserValidateOtpResponse{
			Token: token,
		},
	}

	ctx.JSON(int(statusCode), response)
	return nil
}

func (h *userHandler) GetUserWorkHistoryHandler(ctx echo.Context) error {
	historyCount, history, statusCode, err := h.repo.GetUserWorkHistory(ctx)

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
		Message: "successfully fetched user workhistory",
		Data: map[string]interface{}{
			"total_count": historyCount,
			"history":     history,
		},
	}

	ctx.JSON(int(statusCode), response)
	return nil
}
