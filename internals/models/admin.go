package models

import (
	"io"
	"time"
)

type Admin struct {
	AdminId     string
	Name        string
	Dob         string
	Email       string
	PhoneNumber string
	ProfileUrl  string
	Password    string
	Position    string
}

type CreateAdminRequest struct {
	Name        string `json:"name" validate:"required"`
	Dob         string `json:"dob" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Password    string `json:"password" validate:"required,password"`
	Position    string `json:"position" validate:"required"`
}

type AdminResponse struct {
	AdminId     string    `json:"admin_id"`
	Name        string    `json:"name"`
	Dob         string    `json:"dob"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	ProfileUrl  string    `json:"profile_url"`
	Position    string    `json:"position"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminProfileDetailsResponse struct {
	Name        string    `json:"name"`
	Dob         string    `json:"dob"`
	Email       string    `json:"email"`
	PhoneNumber string    `json:"phone_number"`
	ProfileUrl  string    `json:"profile_url"`
	Position    string    `json:"position"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminForgotPasswordRequest struct {
	Email string `json:"email" validate:"required,email"`
}

type AdminOtpEmailFormat struct {
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	EmailType string            `json:"email_type"`
	Data      map[string]string `json:"data"`
}

type AdminOtpValidateRequest struct {
	Email string `json:"email" validate:"required,email"`
	Otp   string `json:"otp" validate:"required"`
}

type AdminValidateOtpResponse struct {
	Token string `json:"token"`
}

type AdminNewPasswordUpdateRequest struct {
	Password string `json:"password" validate:"required,password"`
}

type AdminProfileUpdateRequest struct {
	AdminId     string `json:"admin_id" validate:"required"`
	Name        string `json:"name" validate:"required"`
	Dob         string `json:"dob" validate:"required"`
	Email       string `json:"email" validate:"required,email"`
	PhoneNumber string `json:"phone_number" validate:"required"`
	Position    string `json:"position" validate:"required"`
}

type AdminInterface interface {
	CheckAdminIdExists(adminId string) (bool, error)
	CreateAdmin(admin *Admin) error
	GetPrevAdminProfileUrl(adminId string) (string, error)
	UpdateAdminProfilePictureUrl(adminId string, url string) error
	StoreAdminOtp(email string, otp string, expireTime time.Time) error
	DeleteAdminOtp(email string, otp string) error
	ValidateAdminOtp(email string, otp string) (bool, error)
	GetAdminDetailsForValidOtp(email string) (string, string, error)
	UpdateAdminNewPassword(adminId string, password string) error
	GetAdminProfileDetails(adminId string) (*AdminProfileDetailsResponse, error)
	GetAllAdmins() ([]*AdminResponse, error)
	CheckAdminEmailsExists(email string) (bool, error)
	GetAdminForLogin(email string) (string, string, string, error)
	UpdateAdminProfileInfo(updateRequest *AdminProfileUpdateRequest) error
}

type AdminStorageInterface interface {
	UploadAdminProfilePicture(fileName string, file io.ReadSeeker) error
	DeleteAdminProfilePicture(fileName string) error
}

type AdminEmailServiceInterface interface {
	SendEmail(data []byte) error
}
