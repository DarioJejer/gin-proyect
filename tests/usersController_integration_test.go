package tests

import (
	"Gin/controllers"
	"Gin/initializers"
	"Gin/models"
	"Gin/repositories"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type IntegrationTestSuite struct {
	suite.Suite
}

func TestIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &IntegrationTestSuite{})
}

var userController *controllers.UsersController

func (suite *IntegrationTestSuite) SetupSuite() {
	// When running go test, the CWD is typically the directory of the test file itself, not the project root where the .env file might be located.
	initializers.LoadEnvVariables("./.env")
	initializers.ConnectToDB()

	userController = controllers.NewUsersController(repositories.NewUsersRepository(), repositories.NewCompaniesRepository())
}

func (suite *IntegrationTestSuite) SetupTest() {

	if respDB := initializers.DB.Exec("DELETE FROM books"); respDB.Error != nil {
		suite.T().Error("Failed to clean up database:", respDB.Error.Error())
	}
	if respDB := initializers.DB.Exec("DELETE FROM users"); respDB.Error != nil {
		suite.T().Error("Failed to clean up database:", respDB.Error.Error())
	}
	if respDB := initializers.DB.Exec("DELETE FROM companies"); respDB.Error != nil {
		suite.T().Error("Failed to clean up database:", respDB.Error.Error())
	}
}

func (suite *IntegrationTestSuite) Test_Post_ValidCreation_StatusCreated() {

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", userController.PostUser)

	r := httptest.NewRecorder()

	company := models.Company{
		Name: "Test Company",
	}
	initializers.DB.Create(&company)

	user := models.User{
		Name:      "Test User",
		Age:       30,
		CompanyID: company.ID,
	}

	marshalledUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
}
func (suite *IntegrationTestSuite) Test_Post_InvalidCompany_StatusBadRequest() {

	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", userController.PostUser)

	r := httptest.NewRecorder()

	user := models.User{
		Name:      "Test User",
		Age:       30,
		CompanyID: 2, // Assuming this company ID does not exist in the database
	}

	marshalledUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *IntegrationTestSuite) Test_Get_Invalidrequest_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", userController.GetUser)

	company := models.Company{
		Name: "Test Company",
	}
	initializers.DB.Create(&company)

	user := models.User{
		Name:      "Test User",
		Age:       30,
		CompanyID: company.ID,
	}

	initializers.DB.Create(&user)

	r := httptest.NewRecorder()

	id := strconv.FormatUint(uint64(user.ID), 10)

	req, _ := http.NewRequest("GET", "/users/"+id, nil)

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
}
