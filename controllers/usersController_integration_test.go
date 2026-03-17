package controllers

import (
	requestDTOs "Gin/dtos/request"
	responseDTOs "Gin/dtos/response"
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

type UsersIntegrationTestSuite struct {
	suite.Suite
}

func TestUsersIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &UsersIntegrationTestSuite{})
}

var userController *UsersController

func (suite *UsersIntegrationTestSuite) SetupSuite() {
	// When running go test, the CWD is typically the directory of the test file itself, not the project root where the .env file might be located.
	initializers.LoadEnvVariables("../.env")
	initializers.ConnectToDB()

	userController = NewUsersController(repositories.NewUsersRepository(), repositories.NewCompaniesRepository())
}

func (suite *UsersIntegrationTestSuite) SetupTest() {

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

func (suite *UsersIntegrationTestSuite) Test_Post_ValidCreation_StatusCreated() {

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
	var resp struct {
		Status string                       `json:"status"`
		User   responseDTOs.UserResponseDTO `json:"user"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "user created", resp.Status)
	assert.Equal(suite.T(), user.Name, resp.User.Name)
	assert.Equal(suite.T(), user.Age, resp.User.Age)
	assert.Equal(suite.T(), user.CompanyID, resp.User.CompanyID)
}

func (suite *UsersIntegrationTestSuite) Test_Post_InvalidUserData_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", userController.PostUser)

	r := httptest.NewRecorder()
	user := requestDTOs.CreateUserDTO{Name: "Test", CompanyID: 1}
	marshalledUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Post_InvalidBody_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/users", userController.PostUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/users", strings.NewReader("invalid {"))
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Post_CompanyNotFound_StatusBadRequest() {

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

func (suite *UsersIntegrationTestSuite) Test_Get_ValidId_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", userController.GetUser)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)
	user := models.User{Name: "Test User", Age: 30, CompanyID: company.ID}
	initializers.DB.Create(&user)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/"+strconv.FormatUint(uint64(user.ID), 10), nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		User responseDTOs.UserResponseDTO `json:"user"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), user.ID, resp.User.ID)
	assert.Equal(suite.T(), user.Name, resp.User.Name)
	assert.Equal(suite.T(), user.Age, resp.User.Age)
	assert.Equal(suite.T(), user.CompanyID, resp.User.CompanyID)
}

func (suite *UsersIntegrationTestSuite) Test_Get_InvalidId_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", userController.GetUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/abc", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Get_UserNotFound_StatusNotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users/:id", userController.GetUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/99999", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_GetUsers_Empty_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", userController.GetUsers)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Users []responseDTOs.UserResponseDTO `json:"users"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), resp.Users, 0)
}

func (suite *UsersIntegrationTestSuite) Test_GetUsers_WithData_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/users", userController.GetUsers)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)
	initializers.DB.Create(&models.User{Name: "User1", Age: 30, CompanyID: company.ID})
	initializers.DB.Create(&models.User{Name: "User2", Age: 25, CompanyID: company.ID})

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Users []responseDTOs.UserResponseDTO `json:"users"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), resp.Users, 2)
	names := []string{resp.Users[0].Name, resp.Users[1].Name}
	assert.Contains(suite.T(), names, "User1")
	assert.Contains(suite.T(), names, "User2")
}

func (suite *UsersIntegrationTestSuite) Test_Update_Valid_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/users/:id", userController.UpdateUser)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)
	user := models.User{Name: "Old", Age: 25, CompanyID: company.ID}
	initializers.DB.Create(&user)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateUserDTO{Name: "Updated", Age: 30, CompanyID: company.ID}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.FormatUint(uint64(user.ID), 10), strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Status string                       `json:"status"`
		User   responseDTOs.UserResponseDTO `json:"user"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "user updated", resp.Status)
	assert.Equal(suite.T(), dto.Name, resp.User.Name)
	assert.Equal(suite.T(), dto.Age, resp.User.Age)
	assert.Equal(suite.T(), dto.CompanyID, resp.User.CompanyID)
	assert.Equal(suite.T(), user.ID, resp.User.ID)
}

func (suite *UsersIntegrationTestSuite) Test_Update_InvalidBody_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/users/:id", userController.UpdateUser)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)
	user := models.User{Name: "User", Age: 30, CompanyID: company.ID}
	initializers.DB.Create(&user)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/"+strconv.FormatUint(uint64(user.ID), 10), strings.NewReader("invalid {"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Update_InvalidUserData_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/users/:id", userController.UpdateUser)

	r := httptest.NewRecorder()
	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)
	user := models.User{Name: "Test", Age: 30, CompanyID: company.ID}
	initializers.DB.Create(&user)
	dto := requestDTOs.CreateUserDTO{Name: "New Name", CompanyID: company.ID}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.FormatUint(uint64(user.ID), 10), strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Update_InvalidId_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/users/:id", userController.UpdateUser)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateUserDTO{Name: "Test", Age: 30, CompanyID: 1}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("PUT", "/users/abc", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Update_UserNotFound_StatusNotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/users/:id", userController.UpdateUser)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateUserDTO{Name: "Test", Age: 30, CompanyID: company.ID}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("PUT", "/users/99999", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *UsersIntegrationTestSuite) Test_Update_InvalidCompany_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.PUT("/users/:id", userController.UpdateUser)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)
	user := models.User{Name: "User", Age: 30, CompanyID: company.ID}
	initializers.DB.Create(&user)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateUserDTO{Name: "Test", Age: 30, CompanyID: 99999}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("PUT", "/users/"+strconv.FormatUint(uint64(user.ID), 10), strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}
