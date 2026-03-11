package controllers

import (
	requestDTOs "Gin/dtos/request"
	responseDTOs "Gin/dtos/response"
	"Gin/models"
	"Gin/repositories"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type CompaniesController struct {
	CompanyRepo repositories.ICompaniesRepository
	validator   *validator.Validate
}

func NewCompaniesController(companyRepo repositories.ICompaniesRepository) *CompaniesController {
	return &CompaniesController{
		CompanyRepo: companyRepo,
		validator:   validator.New(validator.WithRequiredStructEnabled()),
	}
}

// PostCompany godoc
// @Summary Create a new company
// @Description Create a new company with the provided details
// @Tags companies
// @Accept json
// @Produce json
// @Param company body requestDTOs.CreateCompanyDTO true "Company to create"
// @Success 201 {object} responseDTOs.CompanyResponseDTO
// @Failure 400 {object} responseDTOs.ErrorResponseDTO
// @Example 400 {json} BadRequestExample {"code":"INVALID_BODY","message":"Request body is invalid"}
// @Failure 500 {object} responseDTOs.ErrorResponseDTO
// @Example 500 {json} InternalServerErrorExample {"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}
// @Router /companies [post]
func (cc *CompaniesController) PostCompany(ctx *gin.Context) {
	var requestCompany requestDTOs.CreateCompanyDTO

	if err := ctx.ShouldBindJSON(&requestCompany); err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_BODY",
			Message: "Request body is invalid",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	if err := cc.validator.Struct(requestCompany); err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_COMPANY_DATA",
			Message: "Invalid company data",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	modelCompany := models.Company{
		Name: requestCompany.Name,
	}

	if err := cc.CompanyRepo.PostCompany(&modelCompany); err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create company",
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	responseCompany := responseDTOs.CompanyResponseDTO{
		ID:   modelCompany.ID,
		Name: modelCompany.Name,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "company created", "company": responseCompany})
}

// GetCompanies godoc
// @Summary Get all companies
// @Description Retrieve a list of all companies
// @Tags companies
// @Produce json
// @Success 200 {object} []responseDTOs.CompanyResponseDTO
// @Failure 500 {object} responseDTOs.ErrorResponseDTO
// @Example 500 {json} InternalServerErrorExample {"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}
// @Router /companies [get]
func (cc *CompaniesController) GetCompanies(ctx *gin.Context) {
	companies, err := cc.CompanyRepo.GetCompanies()
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to retrieve companies",
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	var responseCompanies []responseDTOs.CompanyResponseDTO
	for _, company := range companies {
		responseCompanies = append(responseCompanies, responseDTOs.CompanyResponseDTO{
			ID:   company.ID,
			Name: company.Name,
		})
	}

	ctx.JSON(http.StatusOK, gin.H{"companies": responseCompanies})
}

// GetCompany godoc
// @Summary Get a company by ID
// @Description Retrieve a company by its ID
// @Tags companies
// @Produce json
// @Param id path int true "Company ID"
// @Success 200 {object} responseDTOs.CompanyResponseDTO
// @Failure 400 {object} responseDTOs.ErrorResponseDTO
// @Example 400 {json} BadRequestExample {"code":"INVALID_COMPANY_ID","message":"Invalid company ID"}
// @Failure 404 {object} responseDTOs.ErrorResponseDTO
// @Example 404 {json} NotFoundExample {"code":"COMPANY_NOT_FOUND","message":"Company not found"}
// @Failure 500 {object} responseDTOs.ErrorResponseDTO "Internal server error"
// @Example 500 {json} InternalServerErrorExample {"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}
// @Router /companies/{id} [get]
func (cc *CompaniesController) GetCompany(ctx *gin.Context) {
	id := ctx.Param("id")
	companyID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_COMPANY_ID",
			Message: "Invalid company ID",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	company, err := cc.CompanyRepo.GetCompany(uint(companyID))
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "COMPANY_NOT_FOUND",
			Message: "Company not found",
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Internal server error",
		})
		return
	}

	responseCompany := responseDTOs.CompanyResponseDTO{
		ID:   company.ID,
		Name: company.Name,
	}

	ctx.JSON(http.StatusOK, gin.H{"company": responseCompany})
}

