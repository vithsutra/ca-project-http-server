package database

import (
	"context"

	"github.com/vithsutra/ca_project_http_server/internals/models"
)

func (repo *PostgresRepo) CreateAdmin(admin *models.Admin) error {
	query := `INSERT INTO admins (
				admin_id,
				name,
				dob,
				email,
				phone_number,
				profile_url,
				password,
				position
			  ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8)`

	_, err := repo.pool.Exec(
		context.Background(),
		query,
		admin.AdminId,
		admin.Name,
		admin.Dob,
		admin.Email,
		admin.PhoneNumber,
		admin.ProfileUrl,
		admin.Password,
		admin.Position,
	)

	return err
}

func (repo *PostgresRepo) GetAllAdmins() ([]*models.AdminResponse, error) {
	query := `SELECT admin_id,name,dob,email,phone_number,profile_url,position,created_at,updated_at FROM admins`

	rows, err := repo.pool.Query(context.Background(), query)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var adminsResponses []*models.AdminResponse

	for rows.Next() {
		var adminResponse models.AdminResponse
		if err := rows.Scan(
			&adminResponse.AdminId,
			&adminResponse.Name,
			&adminResponse.Dob,
			&adminResponse.Email,
			&adminResponse.PhoneNumber,
			&adminResponse.ProfileUrl,
			&adminResponse.Position,
			&adminResponse.CreatedAt,
			&adminResponse.UpdatedAt,
		); err != nil {
			return nil, err
		}

		adminsResponses = append(adminsResponses, &adminResponse)
	}

	return adminsResponses, nil
}

func (repo *PostgresRepo) CheckAdminEmailsExists(email string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM admins WHERE email=$1 )`
	var adminExists bool
	err := repo.pool.QueryRow(context.Background(), query, email).Scan(&adminExists)
	return adminExists, err
}

func (repo *PostgresRepo) GetAdminForLogin(email string) (string, string, string, error) {
	query := `SELECT admin_id,name,password FROM admins WHERE email = $1`
	var adminId string
	var password string
	var name string
	err := repo.pool.QueryRow(context.Background(), query, email).Scan(&adminId, &name, &password)
	return adminId, name, password, err
}
