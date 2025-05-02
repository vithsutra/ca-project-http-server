package models

import (
	"io"
	"time"
)

type User struct {
	UserId      string
	AdminId     string
	Name        string
	Dob         string
	Email       string
	PhoneNumber string
	ProfileUrl  string
	Password    string
	Position    string
	CategoryId  string
}

type CreateUserRequest struct {
	AdminId     string `json:"admin_id" validate:"required"`
	CategoryId  string `json:"category_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Dob         string `json:"dob" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,password"`
	Position    string `json:"position" validate:"required"`
}

type UserProfileDetailsResponse struct {
	Name         string    `json:"name"`
	Dob          string    `json:"dob"`
	Email        string    `json:"email"`
	PhoneNumber  string    `json:"phone_number"`
	ProfileUrl   string    `json:"profile_url"`
	CategoryId   string    `json:"category_id"`
	CategoryName string    `json:"category_name"`
	LoginStatus  bool      `json:"login_status"`
	Latitude     string    `json:"latitude"`
	Longitude    string    `json:"longitude"`
	UpdatedAt    time.Time `json:"updated_at"`
}

type UserResponse struct {
	UserId       string `json:"user_id"`
	Name         string `json:"name"`
	Dob          string `json:"dob"`
	Email        string `json:"email"`
	PhoneNumber  string `json:"phone_number"`
	ProfileUrl   string `json:"profile_url"`
	Position     string `json:"position"`
	CategoryId   string `json:"category_id"`
	CategoryName string `json:"category_name"`
	LoginStatus  bool   `json:"login_status"`
	Latitude     string `json:"latitude"`
	Longitude    string `json:"longitude"`
}

type UserLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type UserLoginResponse struct {
	Token string `json:"token"`
}

type UserWorkLoginRequest struct {
	UserId    string `json:"user_id" validate:"required"`
	LoginDate string `json:"date" validate:"required,date"`
	LoginTime string `json:"time" validate:"required,time"`
	Latitude  string `json:"latitude" validate:"required,latitude"`
	Longitude string `json:"longitude" validate:"required,longitude"`
}

type UserWorkHistory struct {
	UserId       string
	WorkDate     string
	LoginTime    string
	LogoutTime   string
	Latitude     string
	Longitude    string
	UploadedWork string
}

type UserWorkHistoryResponse struct {
	WorkDate     string    `json:"work_date"`
	LoginTime    string    `json:"login_time"`
	LogoutTime   string    `json:"logout_time"`
	Latitude     string    `json:"latitude"`
	Longitude    string    `json:"longitude"`
	UploadedWork string    `json:"uploaded_work"`
	TimeStamp    time.Time `json:"timestamp"`
}

type UserReportPdfDownloadRequest struct {
	UserId    string `query:"user_id" validate:"required"`
	StartDate string `query:"start_date" validate:"required,date"`
	EndDate   string `query:"end_date" validate:"required,date"`
}

type UserWorkHistoryForPdf struct {
	Date        string
	WorkSummary string
	LoginTime   string
	LogoutTime  string
}

type UserReportPdf struct {
	Name     string
	Position string
	History  []*UserWorkHistoryForPdf
}

type UserWorkLogoutRequest struct {
	UserId     string `json:"user_id" validate:"required"`
	LogoutDate string `json:"date" validate:"required,date"`
	LogoutTime string `json:"time" validate:"required,time"`
	Work       string `json:"work" validate:"required"`
}

type UserLeaveRequest struct {
	UserId      string `json:"user_id" validate:"required"`
	LeaveFrom   string `json:"leave_from" validate:"required,date"`
	LeaveTo     string `json:"leave_to" validate:"required,date"`
	LeaveReason string `json:"leave_reason" validate:"required"`
}

type UserLeave struct {
	LeaveId              string
	UserId               string
	LeaveFrom            string
	LeaveTo              string
	LeaveReason          string
	LeaveStatus          string
	LeaveStatusUpdatedBy string
}

type UserLeaveResponse struct {
	LeaveId                string `json:"leave_id"`
	LeaveFrom              string `json:"leave_from"`
	LeaveTo                string `json:"leave_to"`
	LeaveReason            string `json:"leave_reason"`
	LeaveStatus            string `json:"leave_status"`
	LeaveStatusUpdatedBy   string `json:"leave_status_updated_by"`
	LeaveStatusUpdatedDate string `json:"leave_status_updated_date"`
	LeaveStatusUpdtedTime  string `json:"leave_status_updated_time"`
}

type UserPendingLeaveResponse struct {
	UserId         string    `json:"user_id"`
	UserName       string    `json:"user_name"`
	UserEmail      string    `json:"user_email"`
	UserCategory   string    `json:"user_category"`
	LeaveId        string    `json:"leave_id"`
	LeaveFrom      string    `json:"leave_from"`
	LeaveTo        string    `json:"leave_to"`
	LeaveReason    string    `json:"leave_reason"`
	LeaveCreatedAt time.Time `json:"leave_created_at"`
}

type UserProfileInfoUpdateRequest struct {
	UserId      string `json:"user_id" validate:"required"`
	CategoryId  string `json:"category_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Dob         string `json:"dob" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Position    string `json:"position" validate:"required"`
}

type UserLastProfileUpdateTimeResponse struct {
	Time string `json:"time"`
}

type UserNewPasswordUpdateRequest struct {
	Password string `json:"password" validate:"required,password"`
}

type UserForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type UserOtpEmailFormat struct {
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	EmailType string            `json:"email_type"`
	Data      map[string]string `json:"data"`
}

type UserWelcomeEmailFormat struct {
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	EmailType string            `json:"email_type"`
	Data      map[string]string `json:"data"`
}

type UserOtpValidateRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type UserValidateOtpResponse struct {
	Token string `json:"token"`
}

type UserAttendanceSyncRequest struct {
	UserId string `json:"user_id" validate:"required"`
	Date   string `json:"date" validate:"required,date"`
}

type UserDatabaseInterface interface {
	CheckUserEmailExists(email string) (bool, error)
	CreateUser(user *User) error
	GetUserProfileDetails(userId string) (*UserProfileDetailsResponse, error)
	GetUsers(adminId string) ([]*UserResponse, error)
	GetAdminIdByUserId(userId string) (string, error)
	CheckUserIdExists(userId string) (bool, error)
	DeleteUser(userId string) error
	GetUserForLogin(email string) (string, string, string, error)
	CheckUserWorkEntryExists(userId string, date string) (bool, error)
	UserWorkLogin(userWorkHistory *UserWorkHistory) error
	CheckUserWorkLoginEntryExists(userId string, date string) (bool, error)
	UserWorkLogout(userWorkLogoutRequest *UserWorkLogoutRequest) error
	CheckUserPendingLeaveExists(userId string) (bool, error)
	ApplyUserLeave(userLeave *UserLeave) error
	GetAllUsersPendingLeavesCount(adminId string) (int, error)
	GetAllUsersPendingLeaves(adminId string, limit uint32, offset uint32) ([]*UserPendingLeaveResponse, error)
	GetUsersLeavesCount(userId string, leaveStatus string) (int, error)
	GetUserLeaves(userId string, leaveStatus string, limit uint32, offset uint32) ([]*UserLeaveResponse, error)
	CheckLeaveIdExists(leaveId string) (bool, error)
	CheckPendingLeaveExistsByLeaveId(leaveId string) (bool, error)
	CancelUserLeave(leaveId string, userType string) error
	GrantUserLeave(leaveId string) error
	UpdateUserProfileInfo(userId string, userProfileUpdateRequest *UserProfileInfoUpdateRequest) error
	UpdateUserProfileUrl(userId string, url string) error
	GetUserProfileUrl(userId string) (string, error)
	GetUserLastProfileUpdateTime(userId string) (*time.Time, error)
	UpdateNewUserPassword(userId string, password string) error
	StoreUserOtp(email string, otp string, expireTime *time.Time) error
	ClearOtp(email string, otp string) error
	CheckOtpExists(email string, otp string) (bool, error)
	GetUserDetailsForValidateOtp(email string) (string, string, error)
	GetUsersWorkHistoryCount(userId string) (int, error)
	GetUserWorkHistory(userId string, limit uint32, offset uint32) ([]*UserWorkHistoryResponse, error)
	GetUserInfoForPdf(userId string) (string, string, error)
	GetWorkHistoryForPdf(userId, startDate, endDate string) ([]*UserWorkHistoryForPdf, error)
}

type UserStorageInterface interface {
	UploadUserProfilePicture(profilePictureFileName string, file io.ReadSeeker) error
	DeleteUserProfilePicture(fileName string) error
}

type UserEmailServiceInterface interface {
	SendEmail(data []byte) error
}
