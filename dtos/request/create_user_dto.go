package requestDTOs

type CreateUserDTO struct {
	Name      string `json:"name" binding:"required" example:"John Doe" validate:"required"`
	Age       int    `json:"age" binding:"required" example:"30" validate:"required,gt=0"`
	CompanyID uint   `json:"company_id" binding:"required" example:"1" validate:"required,gt=0"`
}
