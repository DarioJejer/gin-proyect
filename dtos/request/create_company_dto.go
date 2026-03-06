package requestDTOs

type CreateCompanyDTO struct {
	Name string `json:"name" binding:"required" example:"Tech Company" validate:"required"`
}
