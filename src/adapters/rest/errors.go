package rest

import (
	"errors"
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	e "github.com/AliceDiNunno/go-nested-traced-error"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
)

var (
	ErrFormValidation     = errors.New("failed to validate form")
	ErrMissingContentType = errors.New("missing content type")
	ErrUploadFileNotFound = errors.New("upload file not found on client side")
	ErrNotFound           = errors.New("endpoint not found")
)

func codeForError(err error) int {
	switch err {
	case ErrFormValidation:
		return http.StatusBadRequest
	case domain.ErrFailedToGetUser:
		return http.StatusUnauthorized
	case domain.ErrAlbumAlreadyExistingWithThisName:
		return http.StatusConflict
	case ErrNotFound, domain.ErrAlbumNotFound, domain.ErrFileNotFound, domain.ErrPictureNotFound:
		return http.StatusNotFound
	case domain.ErrUnableToDeleteObject, domain.ErrUnableToUploadFile, domain.ErrUnableToRetrievePictures:
		return http.StatusInternalServerError
	case ErrUploadFileNotFound:
		return http.StatusUnprocessableEntity
	case domain.ErrUnableToSaveObject:
		return http.StatusInsufficientStorage
	}
	return http.StatusInternalServerError
}

func (rH RoutesHandler) handleError(c *gin.Context, err *e.Error) {
	code := codeForError(err.Err)

	fields := log.Fields{
		"code": code,
		"ip":   c.ClientIP(),
		"path": c.Request.RequestURI,
	}

	authenticatedUser := rH.getAuthenticatedUser(c)

	if authenticatedUser != nil {
		fields["user_id"] = authenticatedUser.UserID
		fields["err"] = &err
	}

	log.WithFields(fields).Error(err.Err.Error())
	c.AbortWithStatusJSON(code, domain.Status{
		Success: false,
		Message: err.Err.Error(),
	})
}

func (rH RoutesHandler) endpointNotFound(c *gin.Context) {
	rH.handleError(c, e.Wrap(ErrNotFound))
}
