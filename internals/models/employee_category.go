package models

type EmployeeCategory struct {
	CategoryId          string
	AdminId             string
	CategoryName        string
	CategoryDescription string
}

type CreateEmployeeCategoryRequest struct {
	AdminId             string `json:"admin_id" validate:"required"`
	CategoryName        string `json:"category_name" validate:"required"`
	CategoryDescription string `json:"category_description" validate:"required"`
}

type EmployeeCategoryResponse struct {
	CategoryId          string `json:"category_id"`
	CategoryName        string `json:"category_name"`
	CategoryDescription string `json:"category_description"`
}

type CategoryInterface interface {
	CheckEmployeeCategoryExists(adminId string, categoryName string) (bool, error)
	CreateEmployeeCategory(category *EmployeeCategory) error
	GetEmployeeCategories(adminId string) ([]*EmployeeCategoryResponse, error)
	CheckEmployeeCategoryIdExists(categoryId string) (bool, error)
	DeleteEmployeeCategory(categoryId string) error
}
