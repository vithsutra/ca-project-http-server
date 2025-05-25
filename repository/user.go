package repository

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/pkg/utils"
)

type UserRepo struct {
	dbRepo           models.UserDatabaseInterface
	storageRepo      models.UserStorageInterface
	emailServiceRepo models.UserEmailServiceInterface
}

func NewUserRepo(
	dbRepo models.UserDatabaseInterface,
	storageRepo models.UserStorageInterface,
	emailServiceRepo models.UserEmailServiceInterface,
) *UserRepo {
	return &UserRepo{
		dbRepo,
		storageRepo,
		emailServiceRepo,
	}
}
func (repo *UserRepo) CreateUser(ctx echo.Context) (string, int32, error) {
	createUserRequest := new(models.CreateUserRequest)

	if err := ctx.Bind(createUserRequest); err != nil {
		return "", 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.RegisterValidation("password", utils.PasswordValidater); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	if err := validation.Struct(createUserRequest); err != nil {
		return "", 400, errors.New("request body validation error")
	}

	userEmailExists, err := repo.dbRepo.CheckUserEmailExists(createUserRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error occurred")
	}

	if userEmailExists {
		return "", 400, errors.New("user email already exists")
	}

	hashedPassword, err := utils.HashPassword(createUserRequest.Password)

	if err != nil {
		log.Println("error occurred while hashing the user password, Error: ", err.Error())
		return "", 500, errors.New("internal server error occurred")
	}

	userId := uuid.NewString()

	user := &models.User{
		UserId:      userId,
		AdminId:     createUserRequest.AdminId,
		CategoryId:  createUserRequest.CategoryId,
		Name:        createUserRequest.Name,
		Dob:         createUserRequest.Dob,
		Email:       createUserRequest.Email,
		PhoneNumber: createUserRequest.PhoneNumber,
		Position:    createUserRequest.Position,
		ProfileUrl:  "pending",
		Password:    hashedPassword,
	}

	if err := repo.dbRepo.CreateUser(user); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	userWelcomeEmailFormat := &models.UserWelcomeEmailFormat{
		To:        createUserRequest.Email,
		EmailType: "welcome",
		Subject:   "Welcome to Vithsutra Technologies",
		Data: map[string]string{
			"user_name":    createUserRequest.Name,
			"service_name": "CA Application",
		},
	}

	jsonBytes, err := json.Marshal(userWelcomeEmailFormat)
	if err != nil {
		log.Println("error occurred while encoding the json, Error:", err.Error())
		return "", 500, errors.New("internal server error occurred")
	}

	if err := repo.emailServiceRepo.SendEmail(jsonBytes); err != nil {
		log.Println("error occurred while sending the welcome email, Error:", err.Error())
		return "", 500, errors.New("internal server error occurred")
	}

	return userId, 201, nil
}

func (repo *UserRepo) GetUserProfileDetails(ctx echo.Context) (*models.UserProfileDetailsResponse, int32, error) {
	userId := ctx.Param("userId")

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return nil, 400, errors.New("user id not exists")
	}

	details, err := repo.dbRepo.GetUserProfileDetails(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error")
	}

	return details, 200, nil

}

func (repo *UserRepo) GetUsers(ctx echo.Context) ([]*models.UserResponse, int32, error) {
	adminId := ctx.Param("adminId")
	usersResponse, err := repo.dbRepo.GetUsers(adminId)
	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}
	if usersResponse == nil {
		return nil, 404, errors.New("users was empty")
	}
	return usersResponse, 200, nil
}

func (repo *UserRepo) DeleteUser(ctx echo.Context) (int32, error) {
	userId := ctx.Param("userId")
	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return 400, errors.New("user id not exists")
	}

	if err := repo.dbRepo.DeleteUser(userId); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) UserLogin(ctx echo.Context) (*models.UserLoginResponse, int32, error) {
	userLoginRequest := new(models.UserLoginRequest)
	if err := ctx.Bind(userLoginRequest); err != nil {
		return nil, 400, errors.New("invalid json request body")
	}

	validation := validator.New()
	if err := validation.Struct(userLoginRequest); err != nil {
		return nil, 400, errors.New("request body validation error")
	}

	userId, userName, hashedPassword, err := repo.dbRepo.GetUserForLogin(userLoginRequest.Email)
	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if err := utils.CheckPassword(hashedPassword, userLoginRequest.Password); err != nil {
		return nil, 401, errors.New("incorrect password")
	}

	adminId, err := repo.dbRepo.GetAdminIdByUserId(userId)
	if err != nil {
		log.Println("error fetching admin_id: ", err.Error())
		return nil, 500, errors.New("failed to fetch admin ID")
	}

	token, err := utils.GenerateToken(userId, userLoginRequest.Email, userName, adminId, "")
	if err != nil {
		log.Println("error occurred while generating token, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	response := &models.UserLoginResponse{
		Token: token,
	}
	return response, 200, nil
}

func (repo *UserRepo) UserWorkLogin(ctx echo.Context) (int32, error) {
	userWorkLoginRequest := new(models.UserWorkLoginRequest)

	if err := ctx.Bind(userWorkLoginRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.RegisterValidation("date", utils.ValidateDate); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.RegisterValidation("time", utils.ValidateTime); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.RegisterValidation("latitude", utils.ValidateLatitude); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.RegisterValidation("longitude", utils.ValidateLongitude); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.Struct(userWorkLoginRequest); err != nil {
		return 400, errors.New("request body validation error")
	}

	userWorkLoginEntryExists, err := repo.dbRepo.CheckUserWorkEntryExists(userWorkLoginRequest.UserId, userWorkLoginRequest.LoginDate)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if userWorkLoginEntryExists {
		return 400, errors.New("duplicate login entries on the same date not allowed")
	}

	userWorkHistory := &models.UserWorkHistory{
		UserId:       userWorkLoginRequest.UserId,
		WorkDate:     userWorkLoginRequest.LoginDate,
		LoginTime:    userWorkLoginRequest.LoginTime,
		LogoutTime:   "pending",
		Latitude:     userWorkLoginRequest.Latitude,
		Longitude:    userWorkLoginRequest.Longitude,
		UploadedWork: "pending",
	}

	if err := repo.dbRepo.UserWorkLogin(userWorkHistory); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) UserWorkLogout(ctx echo.Context) (int32, error) {
	userWorkLogoutRequest := new(models.UserWorkLogoutRequest)

	if err := ctx.Bind(userWorkLogoutRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.RegisterValidation("date", utils.ValidateDate); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.RegisterValidation("time", utils.ValidateTime); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.Struct(userWorkLogoutRequest); err != nil {
		return 400, errors.New("request body validation error")
	}

	userWorkLoginEntryExists, err := repo.dbRepo.CheckUserWorkLoginEntryExists(userWorkLogoutRequest.UserId, userWorkLogoutRequest.LogoutDate)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !userWorkLoginEntryExists {
		return 400, errors.New("logout entry not allowed without login entry")
	}

	if err := repo.dbRepo.UserWorkLogout(userWorkLogoutRequest); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) ApplyUserLeave(ctx echo.Context) (int32, error) {
	userLeaveRequest := new(models.UserLeaveRequest)

	if err := ctx.Bind(userLeaveRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.RegisterValidation("date", utils.ValidateDate); err != nil {
		log.Println("Error occurred while registering validation, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.Struct(userLeaveRequest); err != nil {
		return 400, errors.New("request body validation error")
	}

	if err := utils.CompareDates(userLeaveRequest.LeaveFrom, userLeaveRequest.LeaveTo); err != nil {
		return 400, err
	}

	userPendingLeaveExists, err := repo.dbRepo.CheckUserPendingLeaveExists(userLeaveRequest.UserId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if userPendingLeaveExists {
		return 400, errors.New("pending leave already exists")
	}

	userLeave := &models.UserLeave{
		LeaveId:              uuid.NewString(),
		UserId:               userLeaveRequest.UserId,
		LeaveFrom:            userLeaveRequest.LeaveFrom,
		LeaveTo:              userLeaveRequest.LeaveTo,
		LeaveReason:          userLeaveRequest.LeaveReason,
		LeaveStatus:          "pending",
		LeaveStatusUpdatedBy: "user",
	}

	if err := repo.dbRepo.ApplyUserLeave(userLeave); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) GetAllUsersPendingLeaves(ctx echo.Context) (int32, []*models.UserPendingLeaveResponse, int32, error) {
	adminId := ctx.Param("adminId")
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return 0, nil, 400, errors.New("page paramater must be valid number")
	}

	if pageInt <= 0 {
		pageInt = 1 //default page
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		return 0, nil, 400, errors.New("limit parameter must be valid number")
	}

	if limitInt <= 0 {
		limitInt = 10 //default limit
	}

	offset := (pageInt - 1) * limitInt

	pendingLeavesCount, err := repo.dbRepo.GetAllUsersPendingLeavesCount(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error occurred")
	}

	pendingLeaves, err := repo.dbRepo.GetAllUsersPendingLeaves(adminId, uint32(limitInt), uint32(offset))

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error occurred")
	}

	if pendingLeaves == nil {
		return 0, nil, 404, errors.New("users pending leaves was empty")
	}

	return int32(pendingLeavesCount), pendingLeaves, 200, nil
}

func (repo *UserRepo) GetUserLeaves(ctx echo.Context) (int32, []*models.UserLeaveResponse, int32, error) {
	userId := ctx.Param("userId")
	leaveStatus := ctx.QueryParam("status")
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return 0, nil, 400, errors.New("page paramater must be valid number")
	}

	if pageInt <= 0 {
		pageInt = 1 //default page
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		return 0, nil, 400, errors.New("limit parameter must be valid number")
	}

	if limitInt <= 0 {
		limitInt = 10 //default limit
	}

	offset := (pageInt - 1) * limitInt

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return 0, nil, 400, errors.New("user id not exists")
	}

	if leaveStatus != "" && leaveStatus != "pending" && leaveStatus != "granted" && leaveStatus != "canceled" {
		return 0, nil, 400, errors.New("invalid leave status")
	}

	usersLeaveCount, err := repo.dbRepo.GetUsersLeavesCount(userId, leaveStatus)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error occurred")
	}

	userLeaveResponses, err := repo.dbRepo.GetUserLeaves(userId, leaveStatus, uint32(limitInt), uint32(offset))

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error occurred")
	}

	if userLeaveResponses == nil {
		return 0, nil, 404, errors.New("user leaves was empty")
	}

	return int32(usersLeaveCount), userLeaveResponses, 200, nil

}

func (repo *UserRepo) CancelUserLeave(ctx echo.Context) (int32, error) {
	userId := ctx.Param("userId")
	leaveId := ctx.Param("leaveId")

	requestUrlPath := ctx.Request().URL.Path

	userType := strings.Split(requestUrlPath, "/")[1]

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return 400, errors.New("user id not exists")
	}

	leaveIdExists, err := repo.dbRepo.CheckLeaveIdExists(leaveId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !leaveIdExists {
		return 400, errors.New("leave id not exists")
	}

	leaveExists, err := repo.dbRepo.CheckPendingLeaveExistsByLeaveId(leaveId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !leaveExists {
		return 400, errors.New("leave was not in pending status")
	}

	if err := repo.dbRepo.CancelUserLeave(leaveId, userType); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) GrantUserLeave(ctx echo.Context) (int32, error) {
	leaveId := ctx.Param("leaveId")

	leaveExists, err := repo.dbRepo.CheckLeaveIdExists(leaveId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !leaveExists {
		return 400, errors.New("leave id not exists")
	}

	pendingLeavExists, err := repo.dbRepo.CheckPendingLeaveExistsByLeaveId(leaveId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !pendingLeavExists {
		return 400, errors.New("leave was not in pending status")
	}

	if err := repo.dbRepo.GrantUserLeave(leaveId); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) UpdateUserProfileInfo(ctx echo.Context) (int32, error) {

	userProfileInfoUpdateRequest := new(models.UserProfileInfoUpdateRequest)

	if err := ctx.Bind(userProfileInfoUpdateRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(userProfileInfoUpdateRequest); err != nil {
		return 400, errors.New("request body validation error")
	}

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userProfileInfoUpdateRequest.UserId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return 400, errors.New("invalid user id")
	}

	if err := repo.dbRepo.UpdateUserProfileInfo(userProfileInfoUpdateRequest.UserId, userProfileInfoUpdateRequest); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) UpdateUserProfilePicture(ctx echo.Context) (int32, error) {

	userId := ctx.Param("userId")

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return 400, errors.New("user id not exists")
	}

	file, err := ctx.FormFile("profile_picture")

	if file == nil {
		return 400, errors.New("input file was empty")
	}

	fileNameArr := strings.Split(file.Filename, ".")

	inputFileType := fileNameArr[len(fileNameArr)-1]

	if err != nil {
		return 400, errors.New("unable to get the uploaded image")
	}

	src, err := file.Open()

	if err != nil {
		return 400, errors.New("error occurred while opening the file")
	}

	defer src.Close()

	profilePictureFileName := fmt.Sprintf("%v.%v", userId, inputFileType)

	prevProfileUrl, err := repo.dbRepo.GetUserProfileUrl(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	//skip if profile prev profile picture was not there
	if prevProfileUrl != "pending" {
		prevFileNameArr := strings.Split(prevProfileUrl, ".")

		prevFileType := prevFileNameArr[len(prevFileNameArr)-1]

		prevProfilePictureFileName := fmt.Sprintf("%v.%v", userId, prevFileType)

		if err := repo.storageRepo.DeleteUserProfilePicture(prevProfilePictureFileName); err != nil {
			log.Println("error occurred with aws s3, Error: ", err.Error())
			return 500, errors.New("internal server error occurred")
		}
	}

	rootS3ObjectUrl := os.Getenv("AWS_S3_OBJECT_ROOT_URL")

	if rootS3ObjectUrl == "" {
		log.Println("missing AWS_S3_OBJECT_ROOT_URL env variable")
		return 500, errors.New("internal server error occurred")
	}

	if err := repo.storageRepo.UploadUserProfilePicture(profilePictureFileName, src); err != nil {
		log.Println("error occurred with aws s3, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	profilePictureFileUrl := fmt.Sprintf("%v/users/%v", rootS3ObjectUrl, profilePictureFileName)

	if err := repo.dbRepo.UpdateUserProfileUrl(userId, profilePictureFileUrl); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *UserRepo) DeleteUserProfilePicture(ctx echo.Context) (int32, error) {
	userId := ctx.Param("userId")

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !userIdExists {
		return 400, errors.New("user id not exists")
	}

	prevProfileUrl, err := repo.dbRepo.GetUserProfileUrl(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if prevProfileUrl == "pending" {
		return 400, errors.New("profile picture was not there")
	}

	stringsArr := strings.Split(prevProfileUrl, ".")

	prevFileType := stringsArr[len(stringsArr)-1]

	prevFileName := fmt.Sprintf("%v.%v", userId, prevFileType)

	if err := repo.storageRepo.DeleteUserProfilePicture(prevFileName); err != nil {
		log.Println("error occurred with aws s3, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if err := repo.dbRepo.UpdateUserProfileUrl(userId, "pending"); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil

}

func (repo *UserRepo) GetUserLastProfileUpdateTime(ctx echo.Context) (string, int32, error) {
	userId := ctx.Param("userId")
	lastUpdateTime, err := repo.dbRepo.GetUserLastProfileUpdateTime(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error occurred")
	}

	return lastUpdateTime.String(), 200, nil
}

func (repo *UserRepo) UpdateUserNewPassword(ctx echo.Context) (int32, error) {
	userPasswordUpdateRequest := new(models.UserNewPasswordUpdateRequest)

	if err := ctx.Bind(userPasswordUpdateRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.RegisterValidation("password", utils.PasswordValidater); err != nil {
		log.Println("error occurred while registering the password validator, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := validation.Struct(userPasswordUpdateRequest); err != nil {
		return 400, errors.New("invalid request body format")
	}

	userId := ctx.Param("userId")

	userIdExists, err := repo.dbRepo.CheckUserIdExists(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if !userIdExists {
		return 400, errors.New("user id not exists")
	}

	hashedPassword, err := utils.HashPassword(userPasswordUpdateRequest.Password)

	if err != nil {
		log.Println("error occurred while hashing the password, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if err := repo.dbRepo.UpdateNewUserPassword(userId, hashedPassword); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}
	return 200, nil
}

func (repo *UserRepo) UserForgotPassword(ctx echo.Context) (int32, error) {
	userForgotPasswordRequest := new(models.UserForgotPasswordRequest)

	if err := ctx.Bind(userForgotPasswordRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(userForgotPasswordRequest); err != nil {
		return 400, errors.New("invalid request body format")
	}

	emailsExists, err := repo.dbRepo.CheckUserEmailExists(userForgotPasswordRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error")
	}

	if !emailsExists {
		return 400, errors.New("user email not exists")
	}

	otp, err := utils.GenerateOTP()

	if err != nil {
		log.Println("error occurred while generating the otp,Error: ", err.Error())
		return 500, errors.New("internal server error was occurred")
	}

	expireTime := time.Now().Add(5 * time.Minute)

	if err := repo.dbRepo.StoreUserOtp(userForgotPasswordRequest.Email, otp, &expireTime); err != nil {
		log.Println("error occurred while generating the otp,Error: ", err.Error())
		return 500, errors.New("internal server error was occurred")
	}

	otpMessage := new(models.UserOtpEmailFormat)

	otpMessage.To = userForgotPasswordRequest.Email
	otpMessage.Subject = "Verification Code to Reset Password"
	otpMessage.EmailType = "otp"
	otpMessage.Data = map[string]string{
		"otp":         otp,
		"expire_time": "5",
	}

	jsonBytes, err := json.Marshal(otpMessage)

	if err != nil {
		log.Println("error occurred while encoding the otp message to json, Error: ", err.Error())
		return 500, errors.New("internal server error was occurred")
	}

	if err := repo.emailServiceRepo.SendEmail(jsonBytes); err != nil {
		log.Println("Error occurred while sending the email, Error: ", err.Error())
		return 500, errors.New("internal server error was occurred")
	}

	//clear otp after 5 minute
	go func() {
		time.Sleep(time.Minute * 5)
		if err := repo.dbRepo.ClearOtp(userForgotPasswordRequest.Email, otp); err != nil {
			log.Println("Error occurred while clearing the user otp, Error: ", err.Error())
			return
		}
	}()

	return 200, nil
}

func (user *UserRepo) ValidateUserOtp(ctx echo.Context) (string, int32, error) {
	otpValidateRequest := new(models.UserOtpValidateRequest)

	if err := ctx.Bind(otpValidateRequest); err != nil {
		return "", 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(otpValidateRequest); err != nil {
		return "", 400, errors.New("invalid request body format")
	}

	userEmailExists, err := user.dbRepo.CheckUserEmailExists(otpValidateRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	if !userEmailExists {
		return "", 400, errors.New("email id not exists")
	}

	otpExists, err := user.dbRepo.CheckOtpExists(otpValidateRequest.Email, otpValidateRequest.Otp)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	if !otpExists {
		return "", 400, errors.New("invalid otp")
	}

	if err := user.dbRepo.ClearOtp(otpValidateRequest.Email, otpValidateRequest.Otp); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	userId, userName, err := user.dbRepo.GetUserDetailsForValidateOtp(otpValidateRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	token, err := utils.GenerateToken(userId, otpValidateRequest.Email, userName, "", "")

	if err != nil {
		log.Println("error occurred while generating jwt token, Error: ", err.Error())
		return "", 500, errors.New("internal server error")
	}

	return token, 200, nil
}

func (user *UserRepo) GetUserWorkHistory(ctx echo.Context) (int32, []*models.UserWorkHistoryResponse, int32, error) {
	userId := ctx.Param("userId")
	page := ctx.QueryParam("page")
	limit := ctx.QueryParam("limit")

	if page == "" {
		page = "1"
	}

	if limit == "" {
		limit = "10"
	}

	pageInt, err := strconv.Atoi(page)

	if err != nil {
		return 0, nil, 400, errors.New("page paramater must be valid number")
	}

	if pageInt <= 0 {
		pageInt = 1 //default page
	}

	limitInt, err := strconv.Atoi(limit)

	if err != nil {
		return 0, nil, 400, errors.New("limit parameter must be valid number")
	}

	if limitInt <= 0 {
		limitInt = 10 //default limit
	}

	offset := (pageInt - 1) * limitInt

	workHistoryCount, err := user.dbRepo.GetUsersWorkHistoryCount(userId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error")
	}

	workHistory, err := user.dbRepo.GetUserWorkHistory(userId, uint32(limitInt), uint32(offset))

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 0, nil, 500, errors.New("internal server error")
	}

	if workHistory == nil {
		return 0, nil, 404, errors.New("empty user work history")
	}

	return int32(workHistoryCount), workHistory, 200, nil

}
func (user *UserRepo) GetAllUsersWorkHistory(ctx echo.Context) (int32, []*models.UserWorkHistoryResponse, int32, error) {
	adminId := ctx.Get("admin_id")
	if adminId == nil {
		return 0, nil, 0, echo.NewHTTPError(http.StatusUnauthorized, "admin ID missing in context")
	}

	// Extract query params
	limitQuery := ctx.QueryParam("limit")
	offsetQuery := ctx.QueryParam("offset")

	limit := uint32(10)
	offset := uint32(0)

	if l, err := strconv.Atoi(limitQuery); err == nil {
		limit = uint32(l)
	}
	if o, err := strconv.Atoi(offsetQuery); err == nil {
		offset = uint32(o)
	}

	// Get paginated data
	workHistory, err := user.dbRepo.GetAllUsersWorkHistory(adminId.(string), limit, offset)
	if err != nil {
		return 0, nil, 0, err
	}

	// Get total count from DB
	totalCount, err := user.dbRepo.CountUsersWorkHistory(adminId.(string))
	if err != nil {
		return 0, nil, 0, err
	}

	return int32(len(workHistory)), workHistory, int32(totalCount), nil
}

func (user *UserRepo) GetAllUsersWorkHistoryByAdminId(adminId string, limit, offset uint32) (int32, []*models.UserWorkHistoryResponse, int32, error) {
	workHistory, err := user.dbRepo.GetAllUsersWorkHistory(adminId, limit, offset)
	if err != nil {
		return 0, nil, 0, err
	}

	totalCount, err := user.dbRepo.CountUsersWorkHistory(adminId)
	if err != nil {
		return 0, nil, 0, err
	}

	return int32(len(workHistory)), workHistory, int32(totalCount), nil
}
func (user *UserRepo) DownloadUserWorkHistory(ctx echo.Context) (*models.UserReportPdf, int32, error) {
	userRequest := new(models.UserReportPdfDownloadRequest)

	if err := ctx.Bind(userRequest); err != nil {
		return nil, 400, errors.New("missing or invalid query parameters")
	}

	validate := validator.New()

	if err := validate.RegisterValidation("date", utils.ValidateDate); err != nil {
		log.Println("error occurred while registering the date validation, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if err := validate.Struct(userRequest); err != nil {
		return nil, 400, errors.New("invalid query parameters")
	}

	if err := utils.CompareDates(userRequest.StartDate, userRequest.EndDate); err != nil {
		return nil, 400, errors.New("end date should be greater than start date")
	}

	userName, userCategory, err := user.dbRepo.GetUserInfoForPdf(userRequest.UserId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error")
	}

	history, err := user.dbRepo.GetWorkHistoryForPdf(userRequest.UserId, userRequest.StartDate, userRequest.EndDate)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error")
	}

	userReportPdf := models.UserReportPdf{
		Name:     userName,
		Position: userCategory,
		History:  history,
	}

	return &userReportPdf, 200, nil

}
