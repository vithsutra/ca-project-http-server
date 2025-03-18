package repository

import (
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/models"
	"github.com/vithsutra/ca_project_http_server/pkg/utils"
)

type AdminRepo struct {
	dbRepo          models.AdminInterface
	storageRepo     models.AdminStorageInterface
	emailServieRepo models.AdminEmailServiceInterface
}

func NewAdminRepo(dbRepo models.AdminInterface, storageRepo models.AdminStorageInterface, emailServiceRepo models.AdminEmailServiceInterface) *AdminRepo {
	return &AdminRepo{
		dbRepo:          dbRepo,
		storageRepo:     storageRepo,
		emailServieRepo: emailServiceRepo,
	}
}

func (repo *AdminRepo) AdminLogin(ctx echo.Context) (*models.AdminLoginResponse, int32, error) {
	adminLoginRequest := new(models.AdminLoginRequest)

	if err := ctx.Bind(adminLoginRequest); err != nil {
		return nil, 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(adminLoginRequest); err != nil {
		return nil, 400, errors.New("request body validation error")
	}

	adminEmailsExists, err := repo.dbRepo.CheckAdminEmailsExists(adminLoginRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if !adminEmailsExists {
		return nil, 400, errors.New("admin email not exists")
	}

	adminId, userName, hashedPassword, err := repo.dbRepo.GetAdminForLogin(adminLoginRequest.Email)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if err := utils.CheckPassword(hashedPassword, adminLoginRequest.Password); err != nil {
		return nil, 401, errors.New("incorrect password")
	}

	token, err := utils.GenerateToken(adminId, adminLoginRequest.Email, userName)

	if err != nil {
		log.Println("error occurred while generating the token, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	return &models.AdminLoginResponse{
		Token: token,
	}, 200, nil

}

func (repo *AdminRepo) GetAdminProfileDetails(ctx echo.Context) (*models.AdminProfileDetailsResponse, int32, error) {
	adminId := ctx.Param("adminId")

	adminIdExists, err := repo.dbRepo.CheckAdminIdExists(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}

	if !adminIdExists {
		return nil, 400, errors.New("admin id not exists")
	}

	details, err := repo.dbRepo.GetAdminProfileDetails(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return nil, 500, errors.New("internal server error occurred")
	}
	return details, 200, nil
}

func (repo *AdminRepo) UpdateAdminNewPassword(ctx echo.Context) (int32, error) {

	newPasswordUpdateRequest := new(models.AdminNewPasswordUpdateRequest)

	if err := ctx.Bind(newPasswordUpdateRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	validation.RegisterValidation("password", utils.PasswordValidater)

	if err := validation.Struct(newPasswordUpdateRequest); err != nil {
		return 400, errors.New("invalid request body format")
	}

	adminId := ctx.Param("adminId")

	adminIdExists, err := repo.dbRepo.CheckAdminIdExists(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !adminIdExists {
		return 400, errors.New("admin id not exists")
	}

	hashedPassword, err := utils.HashPassword(newPasswordUpdateRequest.Password)

	if err != nil {
		log.Println("error occurred while hashing the password, Error: ", err.Error())
		return 500, errors.New("internal server errors occurred")
	}

	if err := repo.dbRepo.UpdateAdminNewPassword(adminId, hashedPassword); err != nil {
		log.Println("error occurred while hashing the password, Error: ", err.Error())
		return 500, errors.New("internal server errors occurred")
	}

	return 200, nil

}

func (repo *AdminRepo) UpdateAdminProfileInfo(ctx echo.Context) (int32, error) {
	profileInfoUpdateRequest := new(models.AdminProfileUpdateRequest)

	if err := ctx.Bind(profileInfoUpdateRequest); err != nil {
		return 400, errors.New("invalid json request body")
	}

	validation := validator.New()

	if err := validation.Struct(profileInfoUpdateRequest); err != nil {
		return 400, errors.New("invalid request body format")
	}

	adminIdExists, err := repo.dbRepo.CheckAdminIdExists(profileInfoUpdateRequest.AdminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !adminIdExists {
		return 400, errors.New("admin id not exists")
	}

	if err := repo.dbRepo.UpdateAdminProfileInfo(profileInfoUpdateRequest); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}

func (repo *AdminRepo) UpdateProfilePicture(ctx echo.Context) (int32, error) {
	adminId := ctx.Param("adminId")

	adminIdExists, err := repo.dbRepo.CheckAdminIdExists(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !adminIdExists {
		return 400, errors.New("admin id not exists")
	}

	file, err := ctx.FormFile("profile_picture")

	if err != nil {
		log.Println("error occurred while getting the profile image file from the form data, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if file == nil {
		return 400, errors.New("input file was empty")
	}

	fileNameArr := strings.Split(file.Filename, ".")

	inputFileType := fileNameArr[len(fileNameArr)-1]

	src, err := file.Open()

	if err != nil {
		log.Println("error occurred while opening the uploaded file, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	defer src.Close()

	profilePictureFileName := fmt.Sprintf("%v.%v", adminId, inputFileType)

	adminPrevProfileUrl, err := repo.dbRepo.GetPrevAdminProfileUrl(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if adminPrevProfileUrl != "pending" {
		prevFileNameArr := strings.Split(adminPrevProfileUrl, ".")

		prevFileType := prevFileNameArr[len(prevFileNameArr)-1]

		prevProfilePictureFileName := fmt.Sprintf("%v.%v", adminId, prevFileType)

		if err := repo.storageRepo.DeleteAdminProfilePicture(prevProfilePictureFileName); err != nil {
			log.Println("error occurred with aws s3, Error: ", err.Error())
			return 500, errors.New("internal server error occurred")
		}
	}

	rootS3ObjectUrl := os.Getenv("AWS_S3_OBJECT_ROOT_URL")

	if rootS3ObjectUrl == "" {
		log.Println("missing AWS_S3_OBJECT_ROOT_URL env variable")
		return 500, errors.New("internal server error occurred")
	}

	if err := repo.storageRepo.UploadAdminProfilePicture(profilePictureFileName, src); err != nil {
		log.Println("error occurred with aws s3, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	profilePictureFileUrl := fmt.Sprintf("%v/admins/%v", rootS3ObjectUrl, profilePictureFileName)

	if err := repo.dbRepo.UpdateAdminProfilePictureUrl(adminId, profilePictureFileUrl); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}
	return 200, nil
}

func (repo *AdminRepo) DeleteProfilePicture(ctx echo.Context) (int32, error) {
	adminId := ctx.Param("adminId")

	adminIdExists, err := repo.dbRepo.CheckAdminIdExists(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if !adminIdExists {
		return 400, errors.New("admin id not exists")
	}

	prevAdminProfileUrl, err := repo.dbRepo.GetPrevAdminProfileUrl(adminId)

	if err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if prevAdminProfileUrl == "pending" {
		return 400, errors.New("admin profile picture was not there")
	}

	prevFileArr := strings.Split(prevAdminProfileUrl, ".")

	prevFileType := prevFileArr[len(prevFileArr)-1]

	prevFileName := fmt.Sprintf("%v.%v", adminId, prevFileType)

	if err := repo.storageRepo.DeleteAdminProfilePicture(prevFileName); err != nil {
		log.Println("error occurred with aws s3, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	if err := repo.dbRepo.UpdateAdminProfilePictureUrl(adminId, "pending"); err != nil {
		log.Println("error occurred with database, Error: ", err.Error())
		return 500, errors.New("internal server error occurred")
	}

	return 200, nil
}
