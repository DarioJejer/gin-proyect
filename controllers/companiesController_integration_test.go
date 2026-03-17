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

type CompaniesIntegrationTestSuite struct {
	suite.Suite
}

func TestCompaniesIntegrationTestSuite(t *testing.T) {
	suite.Run(t, &CompaniesIntegrationTestSuite{})
}

var companiesController *CompaniesController

func (suite *CompaniesIntegrationTestSuite) SetupSuite() {
	initializers.LoadEnvVariables("../.env")
	initializers.ConnectToDB()
	companiesController = NewCompaniesController(repositories.NewCompaniesRepository())
}

func (suite *CompaniesIntegrationTestSuite) SetupTest() {
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

func (suite *CompaniesIntegrationTestSuite) Test_Post_ValidCreation_StatusCreated() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/companies", companiesController.PostCompany)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateCompanyDTO{Name: "Tech Company"}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusCreated, r.Code)
	var resp struct {
		Status  string                          `json:"status"`
		Company responseDTOs.CompanyResponseDTO `json:"company"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "company created", resp.Status)
	assert.Equal(suite.T(), dto.Name, resp.Company.Name)
}

func (suite *CompaniesIntegrationTestSuite) Test_Post_InvalidCompanyData_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/companies", companiesController.PostCompany)

	r := httptest.NewRecorder()
	dto := requestDTOs.CreateCompanyDTO{Name: ""}
	body, _ := json.Marshal(dto)
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader(string(body)))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *CompaniesIntegrationTestSuite) Test_Post_InvalidBody_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.POST("/companies", companiesController.PostCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/companies", strings.NewReader("invalid {"))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *CompaniesIntegrationTestSuite) Test_Get_ValidId_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/companies/:id", companiesController.GetCompany)

	company := models.Company{Name: "Test Company"}
	initializers.DB.Create(&company)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/"+strconv.FormatUint(uint64(company.ID), 10), nil)
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

func (suite *CompaniesIntegrationTestSuite) Test_Get_InvalidId_StatusBadRequest() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/companies/:id", companiesController.GetCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/abc", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusBadRequest, r.Code)
}

func (suite *CompaniesIntegrationTestSuite) Test_Get_CompanyNotFound_StatusNotFound() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/companies/:id", companiesController.GetCompany)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies/99999", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusNotFound, r.Code)
}

func (suite *CompaniesIntegrationTestSuite) Test_GetCompanies_Empty_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/companies", companiesController.GetCompanies)

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Companies []responseDTOs.CompanyResponseDTO `json:"companies"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), resp.Companies, 0)
}

func (suite *CompaniesIntegrationTestSuite) Test_GetCompanies_WithData_StatusOk() {
	gin.SetMode(gin.TestMode)
	router := gin.Default()
	router.GET("/companies", companiesController.GetCompanies)

	initializers.DB.Create(&models.Company{Name: "Company1"})
	initializers.DB.Create(&models.Company{Name: "Company2"})

	r := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/companies", nil)
	router.ServeHTTP(r, req)

	assert.Equal(suite.T(), http.StatusOK, r.Code)
	var resp struct {
		Companies []responseDTOs.CompanyResponseDTO `json:"companies"`
	}
	err := json.Unmarshal(r.Body.Bytes(), &resp)
	assert.NoError(suite.T(), err)
	assert.Len(suite.T(), resp.Companies, 2)
	names := []string{resp.Companies[0].Name, resp.Companies[1].Name}
	assert.Contains(suite.T(), names, "Company1")
	assert.Contains(suite.T(), names, "Company2")
}
