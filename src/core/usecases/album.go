package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	"time"
)

func (i interactor) CreateAlbum(user *domain.User, request *Request.CreateAlbumRequest) error {
	album := &domain.Album{
		User:        user.UserID,
		CreatedDate: time.Now(),
		Parent:      nil,
		Name:        request.Name,
		Tags:        nil,
	}

	album.Initialize()

	workflowSameName, err := i.albumRepo.FindByName(user.UserID, request.Name)

	if err == nil && workflowSameName != nil {
		//A workflow with the same name already exists for this user
		return domain.ErrAlbumAlreadyExistingWithThisName
	}

	err = i.albumRepo.CreateAlbum(album)

	if err != nil {
		return err
	}

	return nil
}

func (i interactor) GetUserAlbums(user *domain.User) ([]*domain.Album, error) {
	workflows, err := i.albumRepo.FindByUser(user.UserID)

	if err != nil {
		return nil, err
	}

	return workflows, nil
}

func (i interactor) GetAlbumsContent(user *domain.User, id string) ([]*domain.Picture, error) {
	album, err := i.albumRepo.FindById(user.UserID, id)

	if err != nil {
		return nil, domain.ErrAlbumNotFound
	}

	pictures, err := i.pictureRepo.FindPictures(album)

	if err != nil || pictures == nil {
		return nil, domain.ErrUnableToRetrievePictures
	}

	return pictures, nil
}

func (i interactor) DeleteAlbum(user *domain.User, id string) error {
	album, err := i.albumRepo.FindById(user.UserID, id)

	if err != nil || album == nil {
		return domain.ErrAlbumNotFound
	}

	if album.User != user.UserID {
		return domain.ErrAlbumNotFound
	}

	err = i.albumRepo.DeleteAlbum(album)

	if err != nil {
		return domain.ErrUnableToDeleteObject
	}

	return nil
}
