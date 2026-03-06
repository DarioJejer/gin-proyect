package responseDTOs

type BookResponseDTO struct {
	ID       uint   `json:"id" example:"1"`
	Title    string `json:"title" example:"Book Title"`
	AuthorID uint   `json:"author_id" example:"1"`
}
