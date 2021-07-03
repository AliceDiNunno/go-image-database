package rest

import "github.com/gin-gonic/gin"

func (rH RoutesHandler) pingHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"ping": "pong"})
}
