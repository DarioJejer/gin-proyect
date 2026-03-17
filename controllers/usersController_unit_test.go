package controllers

import (
	requestDTOs "Gin/dtos/request"
	responseDTOs "Gin/dtos/response"
	"Gin/mocks"
	"Gin/models"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"gorm.io/gorm"
)

type UsersUnitTestSuite struct {
	suite.Suite
}

func TestUsersUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UsersUnitTestSuite))
}

func (suite *UsersUnitTestSuite) SetupSuite() {
}

func (suite *UsersUnitTestSuite) SetupTest() {

}

func (suite *UsersUnitTestSuite) Test_Post_ValidCreation_StatusCreated() {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	// userRepositoryMock.On("PostUser", mock.Anything).Return(nil)
	userRepositoryMock.EXPECT().PostUser(mock.Anything).Return(nil)

	company := models.Company{
		Model: gorm.Model{ID: 1}, // How to set embedded struct fields
		Name:  "Test Company",
	}
	companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(&company, nil)
	// companyRepositoryMock.On("GetCompany", mock.Anything).Return(nil, nil)

	router.POST("/users", userControllerMocked.PostUser)

	r := httptest.NewRecorder()

	user := requestDTOs.CreateUserDTO{
		Name:      "Test User",
		Age:       30,
		CompanyID: 1,
	}

	marshalledUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
	var resp struct {
		Status string                      `json:"status"`
		User   responseDTOs.UserResponseDTO `json:"user"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "user created", resp.Status)
	assert.Equal(suite.T(), user.Name, resp.User.Name)
	assert.Equal(suite.T(), user.Age, resp.User.Age)
	assert.Equal(suite.T(), user.CompanyID, resp.User.CompanyID)
	// In unit tests the repo is mocked so the controller does not receive a persisted ID; response ID may be 0.
}

func (suite *UsersUnitTestSuite) Test_Post_InvalidUserData_StatusBadRequest() {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	// userRepositoryMock.EXPECT().PostUser(mock.Anything).Return(errors.New("error"))
	// companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(nil, nil)

	router.POST("/users", userControllerMocked.PostUser)

	r := httptest.NewRecorder()

	user := requestDTOs.CreateUserDTO{
		Name:      "Test User",
		CompanyID: 1,
	}

	marshalledUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Post_InvalidBody_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	router.POST("/users", userControllerMocked.PostUser)

	r := httptest.NewRecorder()

	req, _ := http.NewRequest("POST", "/users", strings.NewReader("invalid json {"))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Post_CompanyNotFound_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(nil, errors.New("company not found"))

	router.POST("/users", userControllerMocked.PostUser)

	r := httptest.NewRecorder()

	user := requestDTOs.CreateUserDTO{
		Name:      "Test User",
		Age:       30,
		CompanyID: 1,
	}
	marshalledUser, _ := json.Marshal(user)

	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Post_RepoError_StatusInternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	company := models.Company{Model: gorm.Model{ID: 1}, Name: "Test Company"}
	companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(&company, nil)
	userRepositoryMock.EXPECT().PostUser(mock.Anything).Return(errors.New("db error"))

	router.POST("/users", userControllerMocked.PostUser)

	r := httptest.NewRecorder()
	user := requestDTOs.CreateUserDTO{Name: "Test User", Age: 30, CompanyID: 1}
	marshalledUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("POST", "/users", strings.NewReader(string(marshalledUser)))

	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Get_ValidId_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	user := &models.User{Model: gorm.Model{ID: 1}, Name: "Test User", Age: 30, CompanyID: 1}
	userRepositoryMock.EXPECT().GetUser(uint(1)).Return(user, nil)

	router.GET("/users/:id", userControllerMocked.GetUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/1", nil)
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

func (suite *UsersUnitTestSuite) Test_Get_InvalidId_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	router.GET("/users/:id", userControllerMocked.GetUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/abc", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Get_UserNotFound_StatusNotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	userRepositoryMock.EXPECT().GetUser(uint(999)).Return(nil, gorm.ErrRecordNotFound)

	router.GET("/users/:id", userControllerMocked.GetUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users/999", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *UsersUnitTestSuite) Test_GetUsers_Valid_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	users := []models.User{{Model: gorm.Model{ID: 1}, Name: "User1", Age: 30, CompanyID: 1}}
	userRepositoryMock.EXPECT().GetUsers().Return(users, nil)

	router.GET("/users", userControllerMocked.GetUsers)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Users []responseDTOs.UserResponseDTO `json:"users"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), resp.Users, 1)
	assert.Equal(suite.T(), "User1", resp.Users[0].Name)
	assert.Equal(suite.T(), uint(1), resp.Users[0].ID)
	assert.Equal(suite.T(), 30, resp.Users[0].Age)
	assert.Equal(suite.T(), uint(1), resp.Users[0].CompanyID)
}

func (suite *UsersUnitTestSuite) Test_GetUsers_RepoError_StatusInternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	userRepositoryMock.EXPECT().GetUsers().Return(nil, errors.New("db error"))

	router.GET("/users", userControllerMocked.GetUsers)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/users", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Update_ValidCreation_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	company := models.Company{Model: gorm.Model{ID: 1}, Name: "Test Company"}
	companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(&company, nil)
	userRepositoryMock.EXPECT().GetUser(uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}, Name: "Old", Age: 25, CompanyID: 1}, nil)
	userRepositoryMock.EXPECT().UpdateUser(mock.Anything).Return(nil)

	router.PUT("/users/:id", userControllerMocked.UpdateUser)

	r := httptest.NewRecorder()
	user := requestDTOs.CreateUserDTO{Name: "Updated", Age: 30, CompanyID: 1}
	marshalledUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(string(marshalledUser)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Status string                      `json:"status"`
		User   responseDTOs.UserResponseDTO `json:"user"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "user updated", resp.Status)
	assert.Equal(suite.T(), user.Name, resp.User.Name)
	assert.Equal(suite.T(), user.Age, resp.User.Age)
	assert.Equal(suite.T(), user.CompanyID, resp.User.CompanyID)
	assert.Equal(suite.T(), uint(1), resp.User.ID)
}

func (suite *UsersUnitTestSuite) Test_Update_InvalidUserData_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	userRepositoryMock.EXPECT().GetUser(uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}, Name: "User", Age: 30, CompanyID: 1}, nil)

	router.PUT("/users/:id", userControllerMocked.UpdateUser)

	r := httptest.NewRecorder()
	user := requestDTOs.CreateUserDTO{Name: "Test", CompanyID: 1}
	marshalledUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(string(marshalledUser)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Update_InvalidBody_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	userRepositoryMock.EXPECT().GetUser(uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}, Name: "User", Age: 30, CompanyID: 1}, nil)

	router.PUT("/users/:id", userControllerMocked.UpdateUser)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader("invalid {"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Update_CompanyNotFound_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	userRepositoryMock.EXPECT().GetUser(uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}, Name: "User", Age: 30, CompanyID: 1}, nil)
	companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(nil, errors.New("not found"))

	router.PUT("/users/:id", userControllerMocked.UpdateUser)

	r := httptest.NewRecorder()
	user := requestDTOs.CreateUserDTO{Name: "Test", Age: 30, CompanyID: 999}
	marshalledUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(string(marshalledUser)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *UsersUnitTestSuite) Test_Update_RepoError_StatusInternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := NewUsersController(&userRepositoryMock, &companyRepositoryMock)

	company := models.Company{Model: gorm.Model{ID: 1}, Name: "Test Company"}
	companyRepositoryMock.EXPECT().GetCompany(mock.Anything).Return(&company, nil)
	userRepositoryMock.EXPECT().GetUser(uint(1)).Return(&models.User{Model: gorm.Model{ID: 1}, Name: "User", Age: 30, CompanyID: 1}, nil)
	userRepositoryMock.EXPECT().UpdateUser(mock.Anything).Return(errors.New("db error"))

	router.PUT("/users/:id", userControllerMocked.UpdateUser)

	r := httptest.NewRecorder()
	user := requestDTOs.CreateUserDTO{Name: "Test", Age: 30, CompanyID: 1}
	marshalledUser, _ := json.Marshal(user)
	req, _ := http.NewRequest("PUT", "/users/1", strings.NewReader(string(marshalledUser)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}
