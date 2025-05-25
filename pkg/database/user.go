package database

import (
	"context"
	"time"

	"github.com/vithsutra/ca_project_http_server/internals/models"
)

func (repo *PostgresRepo) CheckUserEmailExists(email string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE email = $1 )`
	var userEmailExists bool
	err := repo.pool.QueryRow(context.Background(), query, email).Scan(&userEmailExists)
	return userEmailExists, err
}

func (repo *PostgresRepo) CreateUser(user *models.User) error {
	query := `INSERT INTO USERS (
				admin_id,
				category_id,
				user_id,
				name,
				dob,
				email,
				phone_number,
				profile_url,
				password,
				position
			 ) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$9,$10)`

	_, err := repo.pool.Exec(
		context.Background(),
		query,
		user.AdminId,
		user.CategoryId,
		user.UserId,
		user.Name,
		user.Dob,
		user.Email,
		user.PhoneNumber,
		user.ProfileUrl,
		user.Password,
		user.Position,
	)
	return err
}

func (repo *PostgresRepo) GetUserProfileDetails(userId string) (*models.UserProfileDetailsResponse, error) {
	query := `SELECT 
				u.name,
				u.dob,
				u.email,
				u.phone_number,
				u.profile_url,
				u.category_id,
				ec.category_name,
				u.login_status,
				u.latitude,
				u.longitude,
				u.updated_at
			FROM users u
			JOIN employee_category ec ON u.category_id=ec.category_id
			WHERE u.user_id=$1
			 `

	var details models.UserProfileDetailsResponse

	err := repo.pool.QueryRow(context.Background(), query, userId).Scan(
		&details.Name,
		&details.Dob,
		&details.Email,
		&details.PhoneNumber,
		&details.ProfileUrl,
		&details.CategoryId,
		&details.CategoryName,
		&details.LoginStatus,
		&details.Latitude,
		&details.Longitude,
		&details.UpdatedAt,
	)

	return &details, err
}

func (repo *PostgresRepo) GetUsers(adminId string) ([]*models.UserResponse, error) {
	query := `SELECT 
				u.user_id,
				u.name,
				u.dob,
				u.email,
				u.phone_number,
				u.profile_url,
				u.position,
				u.category_id,
				u.login_status,
				u.latitude,
				u.longitude,
				ec.category_name
			FROM users u
			JOIN employee_category ec ON u.category_id=ec.category_id
			WHERE u.admin_id=$1
			`
	rows, err := repo.pool.Query(context.Background(), query, adminId)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var usersResponse []*models.UserResponse

	for rows.Next() {
		var userResponse models.UserResponse
		if err := rows.Scan(
			&userResponse.UserId,
			&userResponse.Name,
			&userResponse.Dob,
			&userResponse.Email,
			&userResponse.PhoneNumber,
			&userResponse.ProfileUrl,
			&userResponse.Position,
			&userResponse.CategoryId,
			&userResponse.LoginStatus,
			&userResponse.Latitude,
			&userResponse.Longitude,
			&userResponse.CategoryName,
		); err != nil {
			return nil, err
		}

		usersResponse = append(usersResponse, &userResponse)
	}

	return usersResponse, nil
}

func (repo *PostgresRepo) CheckUserIdExists(userId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users WHERE user_id = $1 )`
	var userIdExists bool
	err := repo.pool.QueryRow(context.Background(), query, userId).Scan(&userIdExists)
	return userIdExists, err
}

func (repo *PostgresRepo) DeleteUser(userId string) error {
	query := `DELETE FROM users WHERE user_id = $1`
	_, err := repo.pool.Exec(context.Background(), query, userId)
	return err
}

func (repo *PostgresRepo) GetUserForLogin(email string) (string, string, string, error) {
	query := `SELECT user_id,name,password FROM users WHERE email=$1`
	var userId string
	var password string
	var name string
	err := repo.pool.QueryRow(context.Background(), query, email).Scan(&userId, &name, &password)
	return userId, name, password, err
}

func (repo *PostgresRepo) GetAdminIdByUserId(userId string) (string, error) {
	query := `SELECT admin_id FROM users WHERE user_id = $1`
	var adminId string
	err := repo.pool.QueryRow(context.Background(), query, userId).Scan(&adminId)
	return adminId, err
}

func (repo *PostgresRepo) CheckUserWorkEntryExists(userId string, date string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users_history WHERE user_id = $1 AND work_date = $2 )`
	var entryExists bool
	err := repo.pool.QueryRow(context.Background(), query, userId, date).Scan(&entryExists)
	return entryExists, err
}

func (repo *PostgresRepo) UserWorkLogin(userWorkHistory *models.UserWorkHistory) error {
	query1 := `INSERT INTO users_history (
					user_id,
					work_date,
					login_time,
					logout_time,
					latitude,
					longitude,
					uploaded_work
				) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	query2 := `UPDATE users SET 
					work_date=$2,
					login_time=$3,
					logout_time=$4,
					login_status=$5,
					latitude=$6,
					longitude=$7,
					uploaded_work=$8
				WHERE user_id=$1
			   `

	dbConn, err := repo.pool.Acquire(context.Background())

	if err != nil {
		return err
	}

	defer dbConn.Release()

	tx, err := dbConn.Begin(context.Background())

	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		context.Background(),
		query1,
		userWorkHistory.UserId,
		userWorkHistory.WorkDate,
		userWorkHistory.LoginTime,
		userWorkHistory.LogoutTime,
		userWorkHistory.Latitude,
		userWorkHistory.Longitude,
		userWorkHistory.UploadedWork,
	); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	if _, err := tx.Exec(
		context.Background(),
		query2,
		userWorkHistory.UserId,
		userWorkHistory.WorkDate,
		userWorkHistory.LoginTime,
		userWorkHistory.LogoutTime,
		true,
		userWorkHistory.Latitude,
		userWorkHistory.Longitude,
		userWorkHistory.UploadedWork,
	); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}

func (repo *PostgresRepo) CheckUserWorkLoginEntryExists(userId string, date string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users_history WHERE user_id=$1 AND work_date=$2 AND logout_time=$3 ) `
	var entryExists bool
	err := repo.pool.QueryRow(context.Background(), query, userId, date, "pending").Scan(&entryExists)
	return entryExists, err
}

func (repo *PostgresRepo) UserWorkLogout(userWorkLogoutRequest *models.UserWorkLogoutRequest) error {
	query1 := `UPDATE users_history  SET logout_time=$4,uploaded_work=$5 WHERE user_id = $1 AND work_date=$2 AND logout_time=$3`
	query2 := `UPDATE users SET logout_time=$2,login_status=$3,uploaded_work=$4 WHERE user_id=$1`

	dbConn, err := repo.pool.Acquire(context.Background())

	if err != nil {
		return err
	}

	defer dbConn.Release()

	tx, err := dbConn.Begin(context.Background())

	if err != nil {
		return err
	}

	if _, err := tx.Exec(
		context.Background(),
		query1,
		userWorkLogoutRequest.UserId,
		userWorkLogoutRequest.LogoutDate,
		"pending",
		userWorkLogoutRequest.LogoutTime,
		userWorkLogoutRequest.Work,
	); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	if _, err := tx.Exec(
		context.Background(),
		query2,
		userWorkLogoutRequest.UserId,
		userWorkLogoutRequest.LogoutTime,
		false,
		userWorkLogoutRequest.Work,
	); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	if err := tx.Commit(context.Background()); err != nil {
		tx.Rollback(context.Background())
		return err
	}

	return nil
}

func (repo *PostgresRepo) CheckUserPendingLeaveExists(userId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users_leave_history WHERE user_id = $1 AND status=$2 )`
	var leaveExists bool
	err := repo.pool.QueryRow(context.Background(), query, userId, "pending").Scan(&leaveExists)
	return leaveExists, err
}

func (repo *PostgresRepo) CheckLeaveIdExists(leaveId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users_leave_history WHERE leave_id = $1 )`
	var leaveExists bool
	err := repo.pool.QueryRow(context.Background(), query, leaveId).Scan(&leaveExists)
	return leaveExists, err
}

func (repo *PostgresRepo) ApplyUserLeave(userLeave *models.UserLeave) error {
	query := `INSERT INTO users_leave_history (
			 	leave_id,
				user_id,
				leave_from,
				leave_to,
				leave_reason,
				status,
				status_updated_by
			 ) VALUES ($1,$2,$3,$4,$5,$6,$7)`

	_, err := repo.pool.Exec(
		context.Background(),
		query,
		userLeave.LeaveId,
		userLeave.UserId,
		userLeave.LeaveFrom,
		userLeave.LeaveTo,
		userLeave.LeaveReason,
		userLeave.LeaveStatus,
		userLeave.LeaveStatusUpdatedBy,
	)

	return err
}

func (repo *PostgresRepo) GetAllUsersPendingLeavesCount(adminId string) (int, error) {
	query := `SELECT 
				COUNT(*) 
			  FROM users u 
			  JOIN users_leave_history uh ON u.user_id=uh.user_id
			  WHERE u.admin_id=$1 AND uh.status='pending'`

	var count int

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		adminId,
	).Scan(&count)

	return count, err
}

func (repo *PostgresRepo) GetAllUsersPendingLeaves(adminId string, limit uint32, offset uint32) ([]*models.UserPendingLeaveResponse, error) {

	query := `SELECT 
				u.user_id,
				u.name,
				u.email,
				uc.category_name,
				uh.leave_id,
				uh.leave_from,
				uh.leave_to,
				uh.leave_reason,
				uh.created_at
			FROM users u
			JOIN users_leave_history uh ON u.user_id=uh.user_id
			JOIN employee_category uc ON u.category_id=uc.category_id
			WHERE u.admin_id=$1 AND uh.status = 'pending' ORDER BY uh.created_at DESC LIMIT $2 OFFSET $3`

	rows, err := repo.pool.Query(
		context.Background(),
		query,
		adminId,
		limit,
		offset,
	)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var pendingLeaves []*models.UserPendingLeaveResponse

	for rows.Next() {
		var pendingLeave models.UserPendingLeaveResponse

		if err := rows.Scan(
			&pendingLeave.UserId,
			&pendingLeave.UserName,
			&pendingLeave.UserEmail,
			&pendingLeave.UserCategory,
			&pendingLeave.LeaveId,
			&pendingLeave.LeaveFrom,
			&pendingLeave.LeaveTo,
			&pendingLeave.LeaveReason,
			&pendingLeave.LeaveCreatedAt,
		); err != nil {
			return nil, err
		}

		pendingLeaves = append(pendingLeaves, &pendingLeave)
	}

	return pendingLeaves, nil

}
func (repo *PostgresRepo) GetAllUsersWorkHistory(adminId string, limit, offset uint32) ([]*models.UserWorkHistoryResponse, error) {
	query := `SELECT
		u.name,
		uh.work_date,
		uh.login_time,
		uh.logout_time,
		uh.latitude,
		uh.longitude,
		uh.uploaded_work,
		uh.created_at
	FROM users u
	JOIN users_history uh ON u.user_id = uh.user_id
	WHERE u.admin_id = $1
	ORDER BY uh.work_date DESC
	LIMIT $2 OFFSET $3`

	rows, err := repo.pool.Query(context.Background(), query, adminId, limit, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var workHistory []*models.UserWorkHistoryResponse

	for rows.Next() {
		var history models.UserWorkHistoryResponse
		if err := rows.Scan(
			&history.Name,
			&history.WorkDate,
			&history.LoginTime,
			&history.LogoutTime,
			&history.Latitude,
			&history.Longitude,
			&history.UploadedWork,
			&history.TimeStamp,
		); err != nil {
			return nil, err
		}
		workHistory = append(workHistory, &history)
	}

	return workHistory, nil
}

func (repo *PostgresRepo) GetUsersLeavesCount(userId string, leaveStatus string) (int, error) {
	query := ""

	if leaveStatus == "pending" {
		query = `SELECT COUNT(*) FROM users_leave_history WHERE user_id=$1 AND status='pending'`
	} else if leaveStatus == "granted" {
		query = `SELECR COUNT(*) FROM users_leave_history WHERE user_id=$1 AND status='granted'`
	} else if leaveStatus == "canceled" {
		query = `SELECT COUNT(*) FROM users_leave_history WHERE user_id=$1 AND status='canceled'`
	} else {
		query = `SELECT COUNT(*) FROM users_leave_history WHERE user_id=$1`
	}

	var count int

	err := repo.pool.QueryRow(context.Background(), query, userId).Scan(&count)

	return count, err
}

func (repo *PostgresRepo) GetUserLeaves(userId string, leaveStatus string, limit uint32, offset uint32) ([]*models.UserLeaveResponse, error) {
	query := ""
	if leaveStatus == "pending" {
		query = `SELECT 
					leave_id,
					leave_from,
					leave_to,
					leave_reason,
					status,
					status_updated_by,
					updated_at
				FROM users_leave_history WHERE user_id=$1 AND status='pending ORDER BY updated_at DESC LIMIT $2 OFFSET $3'
				`
	} else if leaveStatus == "granted" {
		query = `SELECT 
					leave_id,
					leave_from,
					leave_to,
					leave_reason,
					status,
					status_updated_by,
					updated_at
				FROM users_leave_history WHERE user_id=$1 AND status='granted' ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	} else if leaveStatus == "canceled" {
		query = `SELECT 
					leave_id,
					leave_from,
					leave_to,
					leave_reason,
					status,
					status_updated_by,
					updated_at
				FROM users_leave_history WHERE user_id=$1 AND status='canceled' ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	} else {
		query = `SELECT 
					leave_id,
					leave_from,
					leave_to,
					leave_reason,
					status,
					status_updated_by,
					updated_at
				FROM users_leave_history WHERE user_id=$1 ORDER BY updated_at DESC LIMIT $2 OFFSET $3`
	}

	rows, err := repo.pool.Query(context.Background(), query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	defer rows.Close()

	var userLeaveResponses []*models.UserLeaveResponse

	for rows.Next() {
		var statusUpdatedTime time.Time
		var userLeaveResponse models.UserLeaveResponse

		if err := rows.Scan(
			&userLeaveResponse.LeaveId,
			&userLeaveResponse.LeaveFrom,
			&userLeaveResponse.LeaveTo,
			&userLeaveResponse.LeaveReason,
			&userLeaveResponse.LeaveStatus,
			&userLeaveResponse.LeaveStatusUpdatedBy,
			&statusUpdatedTime,
		); err != nil {
			return nil, err
		}

		userLeaveResponse.LeaveStatusUpdatedDate = statusUpdatedTime.Format("2006-01-02")
		userLeaveResponse.LeaveStatusUpdtedTime = statusUpdatedTime.Format("15:04")

		userLeaveResponses = append(userLeaveResponses, &userLeaveResponse)
	}

	return userLeaveResponses, nil
}

func (repo *PostgresRepo) CheckPendingLeaveExistsByLeaveId(leaveId string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM users_leave_history WHERE leave_id=$1 AND status='pending' )`
	var leaveExists bool
	err := repo.pool.QueryRow(context.Background(), query, leaveId).Scan(&leaveExists)
	return leaveExists, err
}

func (repo *PostgresRepo) CancelUserLeave(leaveId string, userType string) error {
	query := `UPDATE users_leave_history SET status='canceled',status_updated_by=$2 WHERE leave_id=$1`
	_, err := repo.pool.Exec(context.Background(), query, leaveId, userType)
	return err
}

func (repo *PostgresRepo) GrantUserLeave(leaveId string) error {
	query := `UPDATE users_leave_history SET status='granted',status_updated_by='admin' WHERE leave_id=$1`
	_, err := repo.pool.Exec(context.Background(), query, leaveId)
	return err
}

func (repo *PostgresRepo) UpdateUserProfileInfo(userId string, userProfileUpdateRequest *models.UserProfileInfoUpdateRequest) error {
	query := `UPDATE users SET 
				category_id=$2,
				name=$3,
				dob=$4,
				email=$5,
				phone_number=$6,
				position=$7
			 WHERE user_id=$1;`

	_, err := repo.pool.Exec(
		context.Background(),
		query,
		userId,
		userProfileUpdateRequest.CategoryId,
		userProfileUpdateRequest.Name,
		userProfileUpdateRequest.Dob,
		userProfileUpdateRequest.Email,
		userProfileUpdateRequest.PhoneNumber,
		userProfileUpdateRequest.Position,
	)

	return err
}

func (repo *PostgresRepo) UpdateUserProfileUrl(userId string, url string) error {
	query := `UPDATE users SET profile_url = $2 WHERE user_id = $1`
	_, err := repo.pool.Exec(context.Background(), query, userId, url)
	return err
}

func (repo *PostgresRepo) GetUserProfileUrl(userId string) (string, error) {
	query := `SELECT profile_url FROM users WHERE user_id=$1`
	var profileUrl string
	err := repo.pool.QueryRow(context.Background(), query, userId).Scan(&profileUrl)
	return profileUrl, err
}

func (repo *PostgresRepo) GetUserLastProfileUpdateTime(userId string) (*time.Time, error) {
	query := `SELECT updated_at FROM users WHERE user_id=$1`
	var lastUpdateTime time.Time
	if err := repo.pool.QueryRow(context.Background(), query, userId).Scan(&lastUpdateTime); err != nil {
		return nil, err
	}
	return &lastUpdateTime, nil
}

func (repo *PostgresRepo) UpdateNewUserPassword(userId string, password string) error {
	query := `UPDATE users SET password = $2 WHERE user_id = $1`
	_, err := repo.pool.Exec(context.Background(), query, userId, password)
	return err
}

func (repo *PostgresRepo) StoreUserOtp(email string, otp string, expireTime *time.Time) error {
	query := `INSERT INTO user_otps (email,otp,expire_time) VALUES ($1,$2,$3)`
	_, err := repo.pool.Exec(context.Background(), query, email, otp, expireTime)
	return err
}

func (repo *PostgresRepo) ClearOtp(email string, otp string) error {
	query := `DELETE FROM user_otps WHERE email=$1 AND otp=$2`
	_, err := repo.pool.Exec(context.Background(), query, email, otp)
	return err
}

func (repo *PostgresRepo) CheckOtpExists(email string, otp string) (bool, error) {
	query := `SELECT EXISTS ( SELECT 1 FROM user_otps WHERE email=$1 AND otp=$2 )`
	var otpExists bool
	err := repo.pool.QueryRow(context.Background(), query, email, otp).Scan(&otpExists)
	return otpExists, err
}

func (repo *PostgresRepo) GetUserDetailsForValidateOtp(email string) (string, string, error) {
	query := `SELECT user_id,name FROM users WHERE email=$1`
	var userId string
	var userName string
	err := repo.pool.QueryRow(context.Background(), query, email).Scan(&userId, &userName)
	return userId, userName, err
}

func (repo *PostgresRepo) GetUsersWorkHistoryCount(userId string) (int, error) {
	query := `SELECT COUNT(*) FROM users_history WHERE user_id=$1`

	var count int
	err := repo.pool.QueryRow(context.Background(), query, userId).Scan(&count)

	return count, err
}
func (repo *PostgresRepo) GetAllUsersWorkHistoryCount(adminId string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM users u
		JOIN users_history uh ON u.user_id = uh.user_id
		WHERE u.admin_id = $1
	`

	var count int
	err := repo.pool.QueryRow(context.Background(), query, adminId).Scan(&count)
	if err != nil {
		return 0, err
	}

	return count, nil
}

func (repo *PostgresRepo) GetUserWorkHistory(userId string, limit uint32, offset uint32) ([]*models.UserWorkHistoryResponse, error) {
	query := `SELECT 
					work_date,
					login_time,
					logout_time,
					latitude,
					longitude,
					uploaded_work,
					created_at
			 FROM users_history WHERE user_id=$1 ORDER BY created_at DESC LIMIT $2 OFFSET $3`

	var usersWorkHistory []*models.UserWorkHistoryResponse

	rows, err := repo.pool.Query(context.Background(), query, userId, limit, offset)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var userWorkHistory models.UserWorkHistoryResponse

		if err := rows.Scan(
			&userWorkHistory.WorkDate,
			&userWorkHistory.LoginTime,
			&userWorkHistory.LogoutTime,
			&userWorkHistory.Latitude,
			&userWorkHistory.Longitude,
			&userWorkHistory.UploadedWork,
			&userWorkHistory.TimeStamp,
		); err != nil {
			return nil, err
		}

		usersWorkHistory = append(usersWorkHistory, &userWorkHistory)
	}

	return usersWorkHistory, nil

}

func (repo *PostgresRepo) GetUserInfoForPdf(userId string) (string, string, error) {
	query := `SELECT 
	            users.name,
				employee_category.category_name
			  FROM 
			  	users
			  JOIN 
			  	employee_category ON users.category_id=employee_category.category_id
			  WHERE users.user_id = $1`

	var userName, userCategory string

	err := repo.pool.QueryRow(
		context.Background(),
		query,
		userId,
	).Scan(
		&userName,
		&userCategory,
	)

	return userName, userCategory, err
}

func (repo *PostgresRepo) GetWorkHistoryForPdf(userId, startDate, endDate string) ([]*models.UserWorkHistoryForPdf, error) {
	query := `SELECT 
				work_date,
				login_time,
				logout_time,
				uploaded_work
			 FROM 
			 	users_history
			 WHERE user_id=$1 AND work_date BETWEEN $2 AND $3;
			`
	rows, err := repo.pool.Query(
		context.Background(),
		query,
		userId,
		startDate,
		endDate,
	)

	if err != nil {
		return nil, err
	}

	var usersWorkHistory []*models.UserWorkHistoryForPdf
	for rows.Next() {
		var userWorkHistory models.UserWorkHistoryForPdf
		if err := rows.Scan(
			&userWorkHistory.Date,
			&userWorkHistory.LoginTime,
			&userWorkHistory.LogoutTime,
			&userWorkHistory.WorkSummary,
		); err != nil {
			return nil, err
		}

		usersWorkHistory = append(usersWorkHistory, &userWorkHistory)
	}

	return usersWorkHistory, nil
}
func (repo *PostgresRepo) CountUsersWorkHistory(adminId string) (int, error) {
	query := `
		SELECT COUNT(*) 
		FROM users u
		JOIN users_history uh ON u.user_id = uh.user_id
		WHERE u.admin_id = $1
	`
	var count int
	err := repo.pool.QueryRow(context.Background(), query, adminId).Scan(&count)
	if err != nil {
		return 0, err
	}
	return count, nil
}
