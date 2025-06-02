package handlers

import (
	"fmt"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"

	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/pkg/utils"
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

func (h *userHandler) GetAllUsersWorkHistory(ctx echo.Context) error {
	adminId := ctx.Param("adminId")
	if adminId == "" {
		return ctx.JSON(http.StatusBadRequest, &models.ErrorResponse{
			Status: "error",
			Error:  "adminId is required",
		})
	}

	// Default pagination settings
	limit := uint32(10)
	offset := uint32(0)

	// Parse query parameters
	if l, err := strconv.Atoi(ctx.QueryParam("limit")); err == nil && l > 0 {
		limit = uint32(l)
	}
	if o, err := strconv.Atoi(ctx.QueryParam("offset")); err == nil && o >= 0 {
		offset = uint32(o)
	}

	// Get data
	workCount, workHistory, totalCount, err := h.repo.GetAllUsersWorkHistoryByAdminId(adminId, limit, offset)
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, &models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		})
	}

	// Return with pagination metadata
	return ctx.JSON(http.StatusOK, &models.SuccessResponse{
		Status:  "success",
		Message: "successfully fetched all users work history",
		Data: map[string]interface{}{
			"count":        workCount,
			"total_count":  totalCount,
			"history":      workHistory,
			"limit":        limit,
			"offset":       offset,
			"current_page": offset/limit + 1,
			"total_pages":  int(math.Ceil(float64(totalCount) / float64(limit))),
		},
	})
}

func (h *userHandler) DownloadUserReportPdf(ctx echo.Context) error {
	userReportData, statusCode, err := h.repo.DownloadUserWorkHistory(ctx)

	if err != nil {
		response := &models.ErrorResponse{
			Status: "error",
			Error:  err.Error(),
		}
		ctx.JSON(int(statusCode), response)
		return err
	}

	pdfId, err := utils.GenerateUserReportPdf(userReportData)

	if err != nil {
		log.Println("error occurred while generating the pdf, Error: ", err)
		response := &models.ErrorResponse{
			Status: "error",
			Error:  "error occurred while generating the pdf",
		}
		ctx.JSON(int(statusCode), response)
		return err
	}

	ctx.File(fmt.Sprintf("./users_cache/%s.pdf", pdfId))

	if err := os.Remove(fmt.Sprintf("./users_cache/%s.pdf", pdfId)); err != nil {
		log.Println("error occurred while removing the pdf file: Error", err.Error())
		return err
	}

	return nil
}
