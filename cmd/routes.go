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
	root.Use(middlewares.RootMiddleware())
	root.POST("/create/admin", rootHandler.CreateAdminHandler)
	root.GET("/get/admins", rootHandler.GetAllAdminsHandler)

	//auth routes
	auth := e.Group("/auth")
	auth.POST("/login/admin", adminHandler.AdminLoginHandler)
	auth.POST("/login/user", userHandler.UserLoginHandler)
	auth.POST("/forgot/password", userHandler.UserForgotPasswordHandler)
	auth.POST("/validate/otp", userHandler.ValidateUserOtpHandler)

	//admin routes
	admin := e.Group("/admin")
	admin.Use(middlewares.JwtMiddleware())
	admin.POST("/create/employee_category", employeeCategoryHandler.CreateEmployeeCategoryHandler)
	admin.GET("/get/employee_categories/:adminId", employeeCategoryHandler.GetEmployeeCategoriesHandler)
	admin.DELETE("/delete/employee_category/:categoryId", employeeCategoryHandler.DeleteEmployeeCategory)
	admin.POST("/create/user", userHandler.CreateUserHandler)
	admin.GET("/get/users/:adminId", userHandler.GetUsers)
	admin.DELETE("/delete/user/:userId", userHandler.DeleteUser)
	admin.GET("/get/user_work_history/:userId", userHandler.GetUserWorkHistoryHandler)
	admin.GET("/get/user_leaves/:userId", userHandler.GetUserLeavesHandler)
	admin.PATCH("/cancel/user_leave/:userId/:leaveId", userHandler.CancelUserLeaveHandler)
	admin.PATCH("/grant/user_leave/:leaveId", userHandler.GrantUserLeaveHandler)

	//user routes
	user := e.Group("/user")
	// user.Use(middlewares.JwtMiddleware())
	user.POST("/work/login", userHandler.UserWorkLoginHandler)
	user.POST("/work/logout", userHandler.UserWorkLogoutHandler)
	user.POST("/apply/leave", userHandler.ApplyUserLeaveHandler)
	user.PATCH("/cancel/user_leave/:userId/:leaveId", userHandler.CancelUserLeaveHandler)
	user.GET("/get/user_leaves/:userId", userHandler.GetUserLeavesHandler)
	user.PUT("/update/profile_info", userHandler.UserProfileInfoUpdateHandler)
	user.PUT("/update/profile_picture/:userId", userHandler.UpdateUserProfilePictureHandler)
	user.PATCH("/delete/profile_picture/:userId", userHandler.DeleteProfilePictureHandler)
	user.GET("/last_profile_update_time/:userId", userHandler.GetUserLastProfileUpdateTimeHandler)
	user.POST("/update/password/:userId", userHandler.UpdateUserNewPaswordHandler)
	user.POST("/validate/otp", userHandler.ValidateUserOtpHandler)
	user.GET("/get/work_history/:userId", userHandler.GetUserWorkHistoryHandler)

	return e
}
