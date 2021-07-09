package Request

type EditPictureRequest struct {
	Tags EditTagRequest `json:"tags"`
}

type SearchAlbumRequest struct {
	Tags []string `binding:"required"`
}
