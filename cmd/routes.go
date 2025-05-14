package main

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/vithsutra/ca_project_http_server/internals/handlers"
	"github.com/vithsutra/ca_project_http_server/internals/middlewares"
	"github.com/vithsutra/ca_project_http_server/repository"
)

func InitHttpRoutes(e *echo.Echo,
	rootRepo *repository.RootRepo,
	adminRepo *repository.AdminRepo,
	employeeCategoryRepo *repository.EmployeeCategoryRepo,
	userRepo *repository.UserRepo,
) *echo.Echo {

	rootHandler := handlers.NewRootHandler(rootRepo)
	adminHandler := handlers.NewAdminHandler(adminRepo)
	employeeCategoryHandler := handlers.NewEmployeeCategoryHandler(employeeCategoryRepo)
	userHandler := handlers.NewUserHandler(userRepo)

	//cors
	e.Use(middlewares.CorsMiddlware())
	e.OPTIONS("/*", func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	})
	//root routes
	root := e.Group("/r")
	// root.Use(middlewares.RootMiddleware())
	root.POST("/create/admin", rootHandler.CreateAdminHandler)
	root.GET("/get/admins", rootHandler.GetAllAdminsHandler)

	//auth routes
	auth := e.Group("/auth")
	auth.POST("/login/admin", adminHandler.AdminLoginHandler)
	auth.POST("/login/user", userHandler.UserLoginHandler)
	auth.POST("/admin/forgot/password", adminHandler.AdminForgotPasswordHandler)
	auth.POST("/admin/validate/otp", adminHandler.AdminValidateOtpHandler)
	auth.POST("/user/forgot/password", userHandler.UserForgotPasswordHandler)
	auth.POST("/user/validate/otp", userHandler.ValidateUserOtpHandler)

	//admin routes
	admin := e.Group("/admin")
	admin.Use(middlewares.JwtMiddleware())
	admin.GET("/get/profile_details/:adminId", adminHandler.GetAdminProfileDetailsHandler)
	admin.PATCH("/update/password/:adminId", adminHandler.UpdateAdminNewPasswordHandler)
	admin.PUT("/update/profile_info", adminHandler.UpdateAdminProfileInfoHandler)
	admin.PUT("/update/profile_picture/:adminId", adminHandler.UpdateAdminProfilePictureHandler)
	admin.DELETE("/delete/profile_picture/:adminId", adminHandler.DeleteAdminProfilePictureHandler)
	admin.POST("/create/employee_category", employeeCategoryHandler.CreateEmployeeCategoryHandler)
	admin.GET("/get/employee_categories/:adminId", employeeCategoryHandler.GetEmployeeCategoriesHandler)
	admin.DELETE("/delete/employee_category/:categoryId", employeeCategoryHandler.DeleteEmployeeCategory)
	admin.POST("/create/user", userHandler.CreateUserHandler)
	admin.GET("/get/users/:adminId", userHandler.GetUsers)
	admin.DELETE("/delete/user/:userId", userHandler.DeleteUser)
	admin.GET("/get/user_work_history/:userId", userHandler.GetUserWorkHistoryHandler)
	admin.GET("/get/all_users_work_history/:adminId", userHandler.GetAllUsersWorkHistory)

	admin.GET("/get/users_pending_leaves/:adminId", userHandler.GetUserPendingLeavesHandler)
	admin.GET("/get/user_leaves/:userId", userHandler.GetUserLeavesHandler)

	admin.PATCH("/cancel/user_leave/:userId/:leaveId", userHandler.CancelUserLeaveHandler)
	admin.PATCH("/grant/user_leave/:leaveId", userHandler.GrantUserLeaveHandler)
	admin.GET("/download/user/report", userHandler.DownloadUserReportPdf)

	//user routes
	user := e.Group("/user")
	// user.Use(middlewares.JwtMiddleware())
	user.GET("/get/profile_details/:userId", userHandler.GetUserProfileDetailsHandler)
	user.POST("/work/login", userHandler.UserWorkLoginHandler)
	user.POST("/work/logout", userHandler.UserWorkLogoutHandler)
	user.GET("/get/work_history/:userId", userHandler.GetUserWorkHistoryHandler)
	user.POST("/apply/leave", userHandler.ApplyUserLeaveHandler)
	user.PATCH("/cancel/leave/:userId/:leaveId", userHandler.CancelUserLeaveHandler)
	user.GET("/get/leaves/:userId", userHandler.GetUserLeavesHandler)
	user.PUT("/update/profile_info", userHandler.UserProfileInfoUpdateHandler)
	user.PUT("/update/profile_picture/:userId", userHandler.UpdateUserProfilePictureHandler)
	user.PATCH("/delete/profile_picture/:userId", userHandler.DeleteProfilePictureHandler)
	user.GET("/last_profile_update_time/:userId", userHandler.GetUserLastProfileUpdateTimeHandler)
	user.POST("/update/password/:userId", userHandler.UpdateUserNewPaswordHandler)
	user.POST("/validate/otp", userHandler.ValidateUserOtpHandler)

	return e
}
