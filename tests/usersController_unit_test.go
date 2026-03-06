package tests

import (
	"Gin/controllers"
	requestDTOs "Gin/dtos/request"
	"Gin/mocks"
	"Gin/models"
	"encoding/json"
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

type UnitTestSuite struct {
	suite.Suite
}

func TestUnitTestSuite(t *testing.T) {
	suite.Run(t, new(UnitTestSuite))
}

func (suite *UnitTestSuite) SetupSuite() {
}

func (suite *UnitTestSuite) SetupTest() {

}

func (suite *UnitTestSuite) Test_Post_ValidCreation_StatusCreated() {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := controllers.NewUsersController(&userRepositoryMock, &companyRepositoryMock)

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
}

func (suite *UnitTestSuite) Test_Post_InvalidUserData_StatusBadRequest() {

	gin.SetMode(gin.TestMode)
	router := gin.Default()

	userRepositoryMock := mocks.IUsersRepository{}
	companyRepositoryMock := mocks.ICompaniesRepository{}
	userControllerMocked := controllers.NewUsersController(&userRepositoryMock, &companyRepositoryMock)

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
