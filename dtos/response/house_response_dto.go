package responseDTOs

type HouseResponseDTO struct {
	ID      uint   `json:"id" example:"1"`
	Address string `json:"address" example:"123 Main St"`
	Owner   uint   `json:"owner" example:"1"`
}
