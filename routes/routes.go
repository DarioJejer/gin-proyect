package routes

import (
	"Gin/controllers"
	"Gin/repositories"
	"net/http"

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

var userController *controllers.UsersController

func InitializeControllers() {
	userController = controllers.NewUsersController(repositories.NewUsersRepository(), repositories.NewCompaniesRepository())
}

func SetupRoutes(r *gin.Engine) {

	r.GET("/docs/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	InitializeControllers()

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	r.GET("/hello/:name", func(c *gin.Context) {
		name := c.Param("name")
		c.String(http.StatusOK, "Hello %s!", name)
	})
	setUpUserRoutes(r)
}

func setUpUserRoutes(r *gin.Engine) {
	userRoutes := r.Group("/users")
	userRoutes.POST("/", userController.PostUser)
	userRoutes.GET("/", userController.GetUsers)
	userRoutes.GET("/:id", userController.GetUser)
}
