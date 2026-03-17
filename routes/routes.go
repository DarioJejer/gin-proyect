package routes

import (
	"Gin/controllers"
	"Gin/repositories"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

var (
	userController      *controllers.UsersController
	companiesController *controllers.CompaniesController
)

func InitializeControllers() {
	companiesRepo := repositories.NewCompaniesRepository()
	userController = controllers.NewUsersController(repositories.NewUsersRepository(), companiesRepo)
	companiesController = controllers.NewCompaniesController(companiesRepo)
}

func SetupRoutes(r *gin.Engine) {

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	InitializeControllers()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	setUpUserRoutes(r)
	setUpCompanyRoutes(r)
}

func setUpUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	userRoutes.POST("/", userController.PostUser)
	userRoutes.GET("/", userController.GetUsers)
	userRoutes.GET("/:id", userController.GetUser)
	userRoutes.PUT("/:id", userController.UpdateUser)
}

func setUpCompanyRoutes(r *gin.Engine) {
	companyRoutes := r.Group("/companies")
	companyRoutes.POST("/", companiesController.PostCompany)
	companyRoutes.GET("/", companiesController.GetCompanies)
	companyRoutes.GET("/:id", companiesController.GetCompany)
}
