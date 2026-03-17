package controllers

import (
	requestDTOs "Gin/dtos/request"
	responseDTOs "Gin/dtos/response"
	"Gin/helpers"
	"Gin/models"
	"Gin/repositories"
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

type UsersController struct {
	UserRepo      repositories.IUsersRepository
	CompaniesRepo repositories.ICompaniesRepository
	validator     *validator.Validate
}

func NewUsersController(userRepo repositories.IUsersRepository, companiesRepo repositories.ICompaniesRepository) *UsersController {
	return &UsersController{
		UserRepo:      userRepo,
		CompaniesRepo: companiesRepo,
		validator:     validator.New(validator.WithRequiredStructEnabled()),
	}
}

// PostUser godoc
// @Summary Create a new user
// @Description Create a new user with the provided details
// @Tags users
// @Accept json
// @Produce json
// @Param user body requestDTOs.CreateUserDTO true "User to create"
// @Success 201 {object} responseDTOs.UserResponseDTO
// @Failure 400 {object} responseDTOs.ErrorResponseDTO
// @Example 400 {json} BadRequestExample {"code":"INVALID_BODY","message":"Request body is invalid"}
// @Failure 500 {object} responseDTOs.ErrorResponseDTO
// @Example 500 {json} InternalServerErrorExample {"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}
// @Router /users [post]
func (uc *UsersController) PostUser(ctx *gin.Context) {

	var requestUser requestDTOs.CreateUserDTO

	if err := ctx.ShouldBindJSON(&requestUser); err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_BODY",
			Message: "Request body is invalid",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err := uc.validator.Struct(requestUser)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_USER_DATA",
			Message: "Invalid user data",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	_, err = uc.CompaniesRepo.GetCompany(requestUser.CompanyID)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "COMPANY_NOT_FOUND",
			Message: "Company not found",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	modelUser := models.User{
		Name:      requestUser.Name,
		Age:       requestUser.Age,
		CompanyID: requestUser.CompanyID,
	}

	err = uc.UserRepo.PostUser(&modelUser)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to create user",
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	responseUser := responseDTOs.UserResponseDTO{
		ID:        modelUser.ID,
		Name:      modelUser.Name,
		Age:       modelUser.Age,
		CompanyID: modelUser.CompanyID,
	}

	ctx.JSON(http.StatusCreated, gin.H{"status": "user created", "user": responseUser})
}

// GetUsers godoc
// @Summary Get all users
// @Description Retrieve a list of all users with their associated company, books, and house details
// @Tags users
// @Produce json
// @Success 200 {object} []responseDTOs.UserResponseDTO
// @Failure 500 {object} responseDTOs.ErrorResponseDTO
// @Example 500 {json} InternalServerErrorExample {"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}
// @Router /users [get]
func (uc *UsersController) GetUsers(ctx *gin.Context) {

	var modelsUsers []models.User
	modelsUsers, err := uc.UserRepo.GetUsers()
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to retrieve users",
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	responseUsers := helpers.MapUsersToResponseDTOs(modelsUsers)

	ctx.JSON(http.StatusOK, gin.H{"users": responseUsers})
}

// GetUser godoc
// @Summary Get a user by ID
// @Description Retrieve a user by their ID, including associated company, books, and house details
// @Tags users
// @Produce json
// @Param id path int true "User ID"
// @Success 200 {object} responseDTOs.UserResponseDTO
// @Failure 400 {object} responseDTOs.ErrorResponseDTO
// @Example 400 {json} BadRequestExample {"code":"INVALID_USER_ID","message":"Invalid user ID"}
// @Failure 404 {object} responseDTOs.ErrorResponseDTO
// @Example 404 {json} NotFoundExample {"code":"USER_NOT_FOUND","message":"User not found"}
// @Failure 500 {object} responseDTOs.ErrorResponseDTO "Internal server error"
// @Example 500 {json} InternalServerErrorExample {"code":"INTERNAL_SERVER_ERROR","message":"Internal server error"}
// @Router /users/{id} [get]
func (uc *UsersController) GetUser(ctx *gin.Context) {

	id := ctx.Param("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_USER_ID",
			Message: "Invalid user ID",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	user, err := uc.UserRepo.GetUser(uint(userID))
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "USER_NOT_FOUND",
			Message: "User not found",
		}
		if errors.Is(err, gorm.ErrRecordNotFound) {
			ctx.JSON(http.StatusNotFound, errorResponse)
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}
	responseUser := helpers.MapUserToResponseDTO(user)

	ctx.JSON(http.StatusOK, gin.H{"user": responseUser})
}

func (uc *UsersController) UpdateUser(ctx *gin.Context) {

	var requestUser requestDTOs.CreateUserDTO

	id := ctx.Param("id")
	userID, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_USER_ID",
			Message: "Invalid user ID",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	_, err = uc.UserRepo.GetUser(uint(userID))
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			errorResponse := responseDTOs.ErrorResponseDTO{
				Code:    "USER_NOT_FOUND",
				Message: "User not found",
			}
			ctx.JSON(http.StatusNotFound, errorResponse)
			return
		}
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to get user",
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	if err := ctx.ShouldBindJSON(&requestUser); err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_BODY",
			Message: "Request body is invalid",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	err = uc.validator.Struct(requestUser)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INVALID_USER_DATA",
			Message: "Invalid user data",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	_, err = uc.CompaniesRepo.GetCompany(requestUser.CompanyID)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "COMPANY_NOT_FOUND",
			Message: "Company not found",
		}
		ctx.JSON(http.StatusBadRequest, errorResponse)
		return
	}

	modelUser := models.User{
		Model:     gorm.Model{ID: uint(userID)},
		Name:      requestUser.Name,
		Age:       requestUser.Age,
		CompanyID: requestUser.CompanyID,
	}

	err = uc.UserRepo.UpdateUser(&modelUser)
	if err != nil {
		errorResponse := responseDTOs.ErrorResponseDTO{
			Code:    "INTERNAL_SERVER_ERROR",
			Message: "Failed to update user",
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse)
		return
	}

	responseUser := helpers.MapUserToResponseDTO(&modelUser)

	ctx.JSON(http.StatusOK, gin.H{"status": "user updated", "user": responseUser})

}
