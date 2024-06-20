package routers

import (
	"jxb-eprocurement/controllers"
	"jxb-eprocurement/handlers"
	"jxb-eprocurement/middlewares"
	apis "jxb-eprocurement/routers/api"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func SetupRouter(db *gorm.DB) *gin.Engine {
	router := gin.Default()

	roleController := controllers.RoleController{DB: db}
	userController := controllers.UserController{DB: db}
	vehicleController := controllers.VehicleController{DB: db}
	vehicleTypeController := controllers.VehicleTypeController{DB: db}
	loanController := controllers.LoanController{DB: db}

	router.POST("/login", handlers.Authenticate)

	roleRoutes := router.Group("/roles")
	roleRoutes.Use(middlewares.AuthMiddleware())
	roleRoutes.Use(middlewares.AdminMiddleware())
	{
		roleRoutes.POST("/", roleController.CreateRole)
		roleRoutes.GET("/", roleController.GetAllRoles)
		roleRoutes.GET("/:id", roleController.GetRole)
		roleRoutes.PUT("/:id", roleController.UpdateRole)
		roleRoutes.DELETE("/:id", roleController.DeleteRole)
	}

	userRoutes := router.Group("/users")
	userRoutes.Use(middlewares.AuthMiddleware())
	userRoutes.Use(middlewares.AdminMiddleware())
	{
		userRoutes.POST("/", userController.CreateUser)
		userRoutes.GET("/", userController.GetAllUsers)
		userRoutes.GET("/:id", userController.GetUser)
		userRoutes.PUT("/:id", userController.UpdateUser)
		userRoutes.DELETE("/:id", userController.DeleteUser)
	}

	vehicleRoutes := router.Group("/vehicles")
	vehicleRoutes.Use(middlewares.AuthMiddleware())
	{
		vehicleRoutes.GET("/", vehicleController.GetAllVehicles)
		vehicleRoutes.GET("/:id", vehicleController.GetVehicle)
		vehicleRoutes.POST("/", middlewares.AdminMiddleware(), vehicleController.CreateVehicle)
		vehicleRoutes.PUT("/:id", middlewares.AdminMiddleware(), vehicleController.UpdateVehicle)
		vehicleRoutes.DELETE("/:id", middlewares.AdminMiddleware(), vehicleController.DeleteVehicle)

		// Sub-group for /vehicles/type
		vehicleTypeRoutes := vehicleRoutes.Group("/types")
		{
			vehicleTypeRoutes.GET("/", vehicleTypeController.GetAllVehicleTypes)
			vehicleTypeRoutes.GET("/:id", vehicleTypeController.GetVehicleType)
			vehicleTypeRoutes.POST("/", middlewares.AdminMiddleware(), vehicleTypeController.CreateVehicleType)
			vehicleTypeRoutes.PUT("/:id", middlewares.AdminMiddleware(), vehicleTypeController.UpdateVehicleType)
			vehicleTypeRoutes.DELETE("/:id", middlewares.AdminMiddleware(), vehicleTypeController.DeleteVehicleType)
		}
	}

	loanRoutes := router.Group("/loans")
	loanRoutes.Use(middlewares.AuthMiddleware())
	{
		loanRoutes.POST("", loanController.CreateLoan)
		loanRoutes.GET("", loanController.GetAllLoans)
		loanRoutes.GET("/:id", loanController.GetLoan)
		loanRoutes.PUT("/:id", loanController.UpdateLoan)
		loanRoutes.DELETE("/:id", loanController.DeleteLoan)
		// loanRoutes.POST("/:id/approve", middlewares.AdminMiddleware(), loanController.ApproveLoan)
		loanRoutes.POST("/:id/return", loanController.ReturnLoan)
	}

	// Initialize route groups for versioning
	apis.InitRoutes(router, db)

	return router
}
