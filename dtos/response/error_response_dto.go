package responseDTOs

type ErrorResponseDTO struct {
	Code    string `json:"code" example:"USER_NOT_FOUND"`
	Message string `json:"message" example:"User not found"`
}
