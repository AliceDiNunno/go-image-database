package rest

import (
	"github.com/gin-gonic/gin"
)

func SetRoutes(r *gin.Engine, routesHandler RoutesHandler) {
	r.Use(routesHandler.fetchingUserMiddleware())
	r.NoRoute(routesHandler.endpointNotFound)

	user := r.Group("/")

	user.GET("", routesHandler.GetUserAlbumsHandler)    //Get all albums
	user.POST("", routesHandler.CreateUserAlbumHandler) //Create an album

	album := user.Group("/:album", routesHandler.fetchingAlbumMiddleware())
	album.GET("", routesHandler.GetAlbumContentHandler)          //Get an album content
	album.POST("/find", routesHandler.SearchAlbumContentHandler) //Find album
	album.POST("", routesHandler.PostToAlbumHandler)             //post a picture to an album
	album.DELETE("", routesHandler.DeleteAlbumHandler)           //Deletes an album
	album.PATCH("", routesHandler.EditAlbumDataHandler)

	picture := album.Group("/:picture", routesHandler.fetchingPictureMiddleware())
	picture.GET("", routesHandler.GetPictureHandler)        //get a picture
	picture.DELETE("", routesHandler.RemovePictureHandler)  //remove a picture
	picture.PATCH("", routesHandler.EditPictureDataHandler) //set picture data
}
