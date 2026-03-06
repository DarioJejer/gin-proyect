package helpers

import (
	responseDTOs "Gin/dtos/response"
	"Gin/models"
)

func MapUserToResponseDTO(modelUser *models.User) responseDTOs.UserResponseDTO {
	responseUser := responseDTOs.UserResponseDTO{
		ID:        modelUser.ID,
		Name:      modelUser.Name,
		Age:       modelUser.Age,
		CompanyID: modelUser.CompanyID,
		Books: func() []responseDTOs.BookResponseDTO {
			var responseBooks []responseDTOs.BookResponseDTO
			for _, book := range modelUser.Books {
				bookDto := responseDTOs.BookResponseDTO{
					ID:       book.ID,
					Title:    book.Title,
					AuthorID: book.Author,
				}
				responseBooks = append(responseBooks, bookDto)
			}
			return responseBooks
		}(),
		House: func() *responseDTOs.HouseResponseDTO {
			if modelUser.House == nil {
				return nil
			}
			return &responseDTOs.HouseResponseDTO{
				ID:      modelUser.House.ID,
				Address: modelUser.House.Address,
				Owner:   modelUser.House.Owner,
			}
		}(),
	}
	return responseUser
}

func MapUsersToResponseDTOs(modelsUsers []models.User) []responseDTOs.UserResponseDTO {
	var responseUsers []responseDTOs.UserResponseDTO

	for _, user := range modelsUsers {
		responseUser := MapUserToResponseDTO(&user)
		responseUsers = append(responseUsers, responseUser)
	}
	return responseUsers
}
