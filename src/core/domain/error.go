package domain

import "errors"

var (
	ErrFailedToGetUser                  = errors.New("failed to fetch user")
	ErrAlbumAlreadyExistingWithThisName = errors.New("an album already exists with this name")
	ErrAlbumNotFound                    = errors.New("album not found")
	ErrUnableToDeleteObject             = errors.New("unable to delete object")
	ErrUnableToUploadFile               = errors.New("unable to upload file")
	ErrFileNotFound                     = errors.New("file not found")
	ErrUnableToSaveObject               = errors.New("unable to save object")
)
