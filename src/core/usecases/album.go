package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	"github.com/sirupsen/logrus"
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

func (i interactor) UpdateAlbum(user *domain.User, albumId string, request Request.EditAlbumRequest) error {
	album, err := i.albumRepo.FindById(user.UserID, albumId)

	if err != nil || album == nil {
		return domain.ErrAlbumNotFound
	}

	//When we want to add a tag, if it doesn't exist we want to create it
	for _, tagToAdd := range request.Tags.Add {
		tag, err := i.tagRepo.FindTag(tagToAdd)

		if err != nil || tag == nil {
			tag = &domain.Tag{
				Name: tagToAdd,
			}
			tag.Initialize()

			err := i.tagRepo.CreateTag(tag)

			if err != nil {
				logrus.Errorln(err)
				continue
			}
		}

		//Check if tag already exists
		exists := false
		for _, tagInArray := range album.Tags {
			if tagInArray.ID == tag.ID {
				exists = true
			}
		}

		if !exists {
			album.Tags = append(album.Tags, tag)
		}
	}

	//Otherwise if we want to remove a tag that didn't exist, we don't want to create one
	for _, tagToRemove := range request.Tags.Remove {
		tag, err := i.tagRepo.FindTag(tagToRemove)

		if err != nil {
			logrus.Errorln(err)
			continue
		}

		//Check if tag already exists
		for index, tagInArray := range album.Tags {
			if tagInArray.ID == tag.ID {
				album.Tags = append(album.Tags[:index], album.Tags[index+1:]...)
			}
		}
	}

	if request.Visibility.Active {
		album.IsPublic = request.Visibility.Public
	}

	if request.Name.Active && album.Name != "" {
		album.Name = request.Name.Name
	}

	err = i.albumRepo.UpdateAlbum(album)

	if err != nil {
		return err
	}

	return nil
}

func (i interactor) SearchAlbumContent(user *domain.User, albumId string, request Request.SearchAlbumRequest) ([]*domain.SearchPictureResult, error) {
	album, err := i.albumRepo.FindById(user.UserID, albumId)

	if err != nil || album == nil {
		return nil, domain.ErrAlbumNotFound
	}

	pictures, err := i.pictureRepo.SearchPictures(album, request.Tags)

	return pictures, err
}
