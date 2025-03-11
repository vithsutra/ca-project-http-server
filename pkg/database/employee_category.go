package database

import (
	"context"

	"github.com/vithsutra/ca_project_http_server/internals/models"
)

func (repo *PostgresRepo) CheckEmployeeCategoryExists(adminId string, categoryName string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM employee_category WHERE admin_id = $1 AND category_name=$2 )`
	var categoryExists bool
	err := repo.pool.QueryRow(context.Background(), query, adminId, categoryName).Scan(&categoryExists)
	return categoryExists, err
}

func (repo *PostgresRepo) CreateEmployeeCategory(category *models.EmployeeCategory) error {
	query := `INSERT INTO employee_category (category_id,admin_id,category_name,category_description) VALUES ($1,$2,$3,$4)`
	_, err := repo.pool.Exec(context.Background(), query, category.CategoryId, category.AdminId, category.CategoryName, category.CategoryDescription)
	return err
}

func (repo *PostgresRepo) GetEmployeeCategories(adminId string) ([]*models.EmployeeCategoryResponse, error) {
	query := `SELECT category_id,category_name,category_description FROM employee_category WHERE admin_id=$1`
	rows, err := repo.pool.Query(context.Background(), query, adminId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var employeeCategoriesResponse []*models.EmployeeCategoryResponse
	for rows.Next() {
		var employeeCategoryResponse models.EmployeeCategoryResponse
		if err := rows.Scan(&employeeCategoryResponse.CategoryId,
			&employeeCategoryResponse.CategoryName,
			&employeeCategoryResponse.CategoryDescription,
		); err != nil {
			return nil, err
		}

		employeeCategoriesResponse = append(employeeCategoriesResponse, &employeeCategoryResponse)
	}

	return employeeCategoriesResponse, nil
}

func (repo *PostgresRepo) CheckEmployeeCategoryIdExists(categoryId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM employee_category WHERE category_id = $1 )`
	var categoryIdExists bool
	err := repo.pool.QueryRow(context.Background(), query, categoryId).Scan(&categoryIdExists)
	return categoryIdExists, err
}

func (repo *PostgresRepo) DeleteEmployeeCategory(categoryId string) error {
	query := `DELETE FROM employee_category WHERE category_id=$1`
	_, err := repo.pool.Exec(context.Background(), query, categoryId)
	return err
}
