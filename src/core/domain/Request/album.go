package Request

type CreateAlbumRequest struct {
	Name string `json:"name" binding:"required"`
	Public bool `json:"public"`
}
