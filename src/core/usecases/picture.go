package usecases

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/sirupsen/logrus"
	"io"
)

func (i interactor) UploadPicture(user *domain.User, albumId string, file io.Reader, contentType string) error {
	album, err := i.albumRepo.FindById(user.UserID, albumId)

	if err != nil || album == nil {
		return domain.ErrAlbumNotFound
	}

	picture := domain.Picture{
		User: user.UserID,
		//Tags:  nil,
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
