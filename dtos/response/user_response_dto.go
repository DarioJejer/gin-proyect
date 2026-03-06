package responseDTOs

type UserResponseDTO struct {
	ID        uint              `json:"id" example:"1"`
	Name      string            `json:"name" example:"John Doe"`
	Age       int               `json:"age" example:"30"`
	CompanyID uint              `json:"company_id" example:"1"`
	Books     []BookResponseDTO `json:"books,omitempty"`
	House     *HouseResponseDTO `json:"house,omitempty"`
}
