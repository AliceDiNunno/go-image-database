package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/AliceDiNunno/go-image-database/src/core/domain/Request"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"io"
)

func (i interactor) FetchPicture(user *domain.User, album *domain.Album, fileId uuid.UUID) (*domain.Picture, error) {
	if user == nil {
		return nil, domain.ErrFailedToGetUser
	}

	picture, err := i.pictureRepo.FindById(user.UserID, album.ID, fileId)

	if err != nil {
		return nil, err
	}

	return picture, nil
}

func (i interactor) UploadPicture(user *domain.User, album *domain.Album, file io.Reader, contentType string) error {
	album, err := i.albumRepo.FindById(user.UserID, album.ID)

	if err != nil || album == nil {
		return domain.ErrAlbumNotFound
	}

	picture := domain.Picture{
		User:  user.UserID,
		Tags:  nil,
		Album: album,
	}
	picture.Initialize()

	err = i.fileStorage.UploadFile(picture.ID.String(), file)

	if err != nil {
		return err
	}

	err = i.pictureRepo.CreatePicture(&picture)

	if err != nil {
		err := i.fileStorage.DeleteFile(picture.ID.String())
		if err != nil { //Internal error, we don't want this one to be passed to the client, so we're just logging it
			logrus.Errorln(err)
		}
		return domain.ErrUnableToSaveObject
	}

	return nil
}

func (i interactor) DeletePicture(user *domain.User, album *domain.Album, picture *domain.Picture) error {
	picture, err := i.pictureRepo.FindById(user.UserID, album.ID, picture.ID)

	if err != nil || picture == nil {
		return domain.ErrPictureNotFound
	}

	err = i.pictureRepo.DeletePicture(picture)
	if err != nil {
		return domain.ErrUnableToDeleteObject
	}

	err = i.fileStorage.DeleteFile(picture.ID.String())
	if err != nil {
		return domain.ErrUnableToDeleteObject
	}

	return nil
}

func (i interactor) FetchPictureData(user *domain.User, album *domain.Album, picture *domain.Picture) (io.Reader, error) {
	io, err := i.fileStorage.GetFile(picture.ID.String())
	if err != nil {
		return nil, err
	}

	return io, nil
}

func (i interactor) UpdatePicture(user *domain.User, album *domain.Album, picture *domain.Picture, request Request.EditPictureRequest) error {
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
		for _, tagInArray := range picture.Tags {
			if tagInArray.ID == tag.ID {
				exists = true
			}
		}

		if !exists {
			picture.Tags = append(picture.Tags, tag)
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
		for index, tagInArray := range picture.Tags {
			if tagInArray.ID == tag.ID {
				picture.Tags = append(picture.Tags[:index], picture.Tags[index+1:]...)
			}
		}
	}

	err := i.pictureRepo.UpdatePicture(picture)

	return err
}
