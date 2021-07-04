package Request

type EditTagRequest struct {
	Add    []string `json:"add"`
	Remove []string `json:"remove"`
}

type EditPictureRequest struct {
	Tags EditTagRequest `json:"tags"`
}
