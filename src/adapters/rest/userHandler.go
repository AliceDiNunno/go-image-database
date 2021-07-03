package rest

import (
	"github.com/AliceDiNunno/go-image-database/src/core/domain"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

//TODO: Currently a temporary dummy function for tests
func (rH RoutesHandler) fetchingUserMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		userFound := true

		if !userFound {
			rH.handleError(c, domain.ErrFailedToGetUser)
			return
		}

		bytes := []byte{0xC4, 0x88, 0x10, 0x53, 0x53, 0x01, 0x4c, 0x57, 0x81, 0xfc, 0xcc, 0xac, 0x77, 0x15, 0x13, 0xcc}

		id, err := uuid.FromBytes(bytes)

		if err != nil {
			rH.handleError(c, err)
			return
		}

		user := domain.User{
			UserID: id,
		}
		c.Set("authenticatedUser", user)
		c.Next()
	}
}

//if fetchingUserMiddleware is a root middleware there is no reason than domain.user is empty at this point
//We will assume fetchingUserMiddleware is a root middleware
func (rH RoutesHandler) getAuthenticatedUser(c *gin.Context) *domain.User {
	auth, exists := c.Get("authenticatedUser")

	if !exists {
		return nil
	}

	authenticatedUser := auth.(domain.User)

	return &authenticatedUser
}
