package database

import (
	"context"
	"time"

	"github.com/vithsutra/ca_project_http_server/internals/models"
)

func (repo *PostgresRepo) CheckAdminIdExists(adminId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM admins WHERE admin_id=$1 )`
	var exists bool
	err := repo.pool.QueryRow(context.Background(), query, adminId).Scan(&exists)
	return exists, err
}

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

func (repo *PostgresRepo) GetAdminProfileDetails(adminId string) (*models.AdminProfileDetailsResponse, error) {
	query := `SELECT 
				name,
				dob,
				email,
				phone_number,
				profile_url,
				position,
				updated_at
			FROM admins WHERE admin_id=$1`

	var details models.AdminProfileDetailsResponse

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		adminId,
	).Scan(
		&details.Name,
		&details.Dob,
		&details.Email,
		&details.PhoneNumber,
		&details.ProfileUrl,
		&details.Position,
		&details.UpdatedAt,
	)

	return &details, err
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

func (repo *PostgresRepo) UpdateAdminNewPassword(adminId string, password string) error {
	query := `UPDATE admins SET password=$2 WHERE admin_id=$1`
	_, err := repo.pool.Exec(context.Background(), query, adminId, password)
	return err
}

func (repo *PostgresRepo) StoreAdminOtp(email string, otp string, expireTime time.Time) error {
	query := `INSERT INTO (email,otp,expire_time) admins_otp VALUES ($1,$2,$3)`
	_, err := repo.pool.Exec(
		context.Background(),
		query,
		email,
		otp,
		expireTime,
	)
	return err
}

func (repo *PostgresRepo) DeleteAdminOtp(email string, otp string) error {
	query := `DELETE FROM admin_otps WHERE email=$1 AND otp=$2`
	_, err := repo.pool.Exec(
		context.Background(),
		query,
		email,
		otp,
	)
	return err
}

func (repo *PostgresRepo) ValidateAdminOtp(email string, otp string) (bool, error) {
	query1 := `SELECT EXISTS ( SELECT 1 FROM admin_otps WHERE email=$1 AND otp=$2 )`
	query2 := `DELETE FROM admin_otps WHERE email=$1 AND otp=$2`

	var otpExists bool

	err := repo.pool.QueryRow(
		context.Background(),
		query1,
		email,
		otp,
	).Scan(&otpExists)

	if err != nil {
		return false, err
	}

	if !otpExists {
		return false, nil
	}

	_, err = repo.pool.Exec(
		context.Background(),
		query2,
		email,
		otp,
	)

	if err != nil {
		return false, err
	}

	return true, nil
}

func (repo *PostgresRepo) GetAdminDetailsForValidOtp(email string) (string, string, error) {
	query := `SELECT admin_id,name FROM admins WHERE email=$1`

	var adminId string
	var adminName string

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		email,
	).Scan(&adminId, &adminName)

	return adminId, adminName, err
}

func (repo *PostgresRepo) UpdateAdminProfileInfo(updateRequest *models.AdminProfileUpdateRequest) error {
	query := `UPDATE admins SET 
				name=$2,
				dob=$3,
				email=$4,
				phone_number=$5,
				position=$6
			WHERE admin_id=$1`

	_, err := repo.pool.Exec(
		context.Background(),
		query,
		updateRequest.AdminId,
		updateRequest.Name,
		updateRequest.Dob,
		updateRequest.Email,
		updateRequest.PhoneNumber,
		updateRequest.Position,
	)

	return err
}

func (repo *PostgresRepo) GetPrevAdminProfileUrl(adminId string) (string, error) {
	query := `SELECT profile_url FROM admins WHERE admin_id=$1`
	var profileUrl string
	err := repo.pool.QueryRow(context.Background(), query, adminId).Scan(&profileUrl)
	return profileUrl, err
}

func (repo PostgresRepo) UpdateAdminProfilePictureUrl(adminId string, url string) error {
	query := `UPDATE admins SET profile_url=$2 WHERE admin_id=$1`
	_, err := repo.pool.Exec(context.Background(), query, adminId, url)
	return err
}
