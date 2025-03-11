package models

import "time"

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

type AdminLoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required"`
}

type AdminLoginResponse struct {
	Token string `json:"token"`
}

type AdminInterface interface {
	CreateAdmin(admin *Admin) error
	GetAllAdmins() ([]*AdminResponse, error)
	CheckAdminEmailsExists(email string) (bool, error)
	GetAdminForLogin(email string) (string, string, string, error)
}
