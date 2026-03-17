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

type CompaniesUnitTestSuite struct {
	suite.Suite
}

func TestCompaniesUnitTestSuite(t *testing.T) {
	suite.Run(t, new(CompaniesUnitTestSuite))
}

func (suite *CompaniesUnitTestSuite) SetupSuite() {}

func (suite *CompaniesUnitTestSuite) SetupTest() {}

func (suite *CompaniesUnitTestSuite) Test_Post_ValidCreation_StatusCreated() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	companyRepositoryMock.EXPECT().PostCompany(mock.Anything).Return(nil)

	router.POST("/companies", companiesControllerMocked.PostCompany)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateCompanyDTO{Name: "Tech Company"}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
	var resp struct {
		Status  string                          `json:"status"`
		Company responseDTOs.CompanyResponseDTO  `json:"company"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "company created", resp.Status)
	assert.Equal(suite.T(), dto.Name, resp.Company.Name)
	// In unit tests the repo is mocked so the controller does not receive a persisted ID; response ID may be 0.
}

func (suite *CompaniesUnitTestSuite) Test_Post_InvalidCompanyData_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	router.POST("/companies", companiesControllerMocked.PostCompany)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateCompanyDTO{Name: ""}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *CompaniesUnitTestSuite) Test_Post_InvalidBody_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	router.POST("/companies", companiesControllerMocked.PostCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader("invalid {"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *CompaniesUnitTestSuite) Test_Post_RepoError_StatusInternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	companyRepositoryMock.EXPECT().PostCompany(mock.Anything).Return(errors.New("db error"))

	router.POST("/companies", companiesControllerMocked.PostCompany)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateCompanyDTO{Name: "Tech Company"}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *CompaniesUnitTestSuite) Test_Get_ValidId_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	company := &models.Company{Model: gorm.Model{ID: 1}, Name: "Tech Company"}
	companyRepositoryMock.EXPECT().GetCompany(uint(1)).Return(company, nil)

	router.GET("/companies/:id", companiesControllerMocked.GetCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/1", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Company responseDTOs.CompanyResponseDTO `json:"company"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), company.ID, resp.Company.ID)
	assert.Equal(suite.T(), company.Name, resp.Company.Name)
}

func (suite *CompaniesUnitTestSuite) Test_Get_InvalidId_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	router.GET("/companies/:id", companiesControllerMocked.GetCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/abc", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *CompaniesUnitTestSuite) Test_Get_CompanyNotFound_StatusNotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	companyRepositoryMock.EXPECT().GetCompany(uint(999)).Return(nil, gorm.ErrRecordNotFound)

	router.GET("/companies/:id", companiesControllerMocked.GetCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/999", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *CompaniesUnitTestSuite) Test_Get_RepoError_StatusInternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	companyRepositoryMock.EXPECT().GetCompany(uint(1)).Return(nil, errors.New("db error"))

	router.GET("/companies/:id", companiesControllerMocked.GetCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/1", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}

func (suite *CompaniesUnitTestSuite) Test_GetCompanies_Valid_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	companies := []models.Company{{Model: gorm.Model{ID: 1}, Name: "Company1"}}
	companyRepositoryMock.EXPECT().GetCompanies().Return(companies, nil)

	router.GET("/companies", companiesControllerMocked.GetCompanies)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Companies []responseDTOs.CompanyResponseDTO `json:"companies"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), resp.Companies, 1)
	assert.Equal(suite.T(), "Company1", resp.Companies[0].Name)
	assert.Equal(suite.T(), uint(1), resp.Companies[0].ID)
}

func (suite *CompaniesUnitTestSuite) Test_GetCompanies_RepoError_StatusInternalServerError() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()

	companyRepositoryMock := mocks.ICompaniesRepository{}
	companiesControllerMocked := NewCompaniesController(&companyRepositoryMock)

	companyRepositoryMock.EXPECT().GetCompanies().Return(nil, errors.New("db error"))

	router.GET("/companies", companiesControllerMocked.GetCompanies)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusInternalServerError, r.Code)
}
