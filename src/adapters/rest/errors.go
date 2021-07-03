package rest

import (
	"errors"
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"net/http"
	"runtime"
	"strings"
)

var (
	ErrFormValidation     = errors.New("failed to validate form")
	ErrUploadFileNotFound = errors.New("upload file not found on client side")
	ErrNotFound           = errors.New("endpoint not found")
)

func getFrame(skipFrames int) runtime.Frame {
	// We need the frame at index skipFrames+2, since we never want runtime.Callers and getFrame
	targetFrameIndex := skipFrames + 2

	// Set size to targetFrameIndex+2 to ensure we have room for one more caller than we need
	programCounters := make([]uintptr, targetFrameIndex+2)
	n := runtime.Callers(0, programCounters)

	frame := runtime.Frame{Function: "unknown"}
	if n > 0 {
		frames := runtime.CallersFrames(programCounters[:n])
		for more, frameIndex := true, 0; more && frameIndex <= targetFrameIndex; frameIndex++ {
			var frameCandidate runtime.Frame
			frameCandidate, more = frames.Next()
			if frameIndex == targetFrameIndex {
				frame = frameCandidate
			}
		}
	}

	return frame
}

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

func getFunctionName(depth int) string {
	function := getFrame(depth).Function
	functionSplitted := strings.Split(function, "/")
	functionName := functionSplitted[len(functionSplitted)-1:][0]

	specifiedFunctionActionSplitted := strings.Split(functionName, ".")
	specifiedFunctionName := specifiedFunctionActionSplitted[2]

	return specifiedFunctionName
}

func (rH RoutesHandler) handleError(c *gin.Context, err error) {
	var depth = 2
	errName := getFunctionName(depth) + ": " + err.Error()
	code := codeForError(err)

	fields := log.Fields{
		"code": code,
		"ip":   c.ClientIP(),
		"path": c.FullPath(),
	}

	authenticatedUser := rH.getAuthenticatedUser(c)

	if authenticatedUser != nil {
		fields["user_id"] = authenticatedUser.UserID
	}

	log.WithFields(fields).Error(errors.New(errName))
	c.AbortWithStatusJSON(code, domain.Status{
		Success: false,
		Message: err.Error(),
	})
}

func (rH RoutesHandler) endpointNotFound(c *gin.Context) {
	rH.handleError(c, ErrNotFound)
}
