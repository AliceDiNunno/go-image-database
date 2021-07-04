package Request

type CreateAlbumRequest struct {
	Name   string `json:"name" binding:"required"`
	Public bool   `json:"public"`
}

type EditAlbumNameRequest struct {
	Active bool   `json:"active"`
	Name   string `json:"name"`
}

type EditAlbumVisibilityRequest struct {
	Active bool `json:"active"`
	Public bool `json:"public"`
}

type EditAlbumRequest struct {
	Tags       EditTagRequest             `json:"tags"`
	Name       EditAlbumNameRequest       `json:"name"`
	Visibility EditAlbumVisibilityRequest `json:"visibility"`
}
