package api

import (
	"github.com/crackeer/go-skeleton/util"

	"github.com/gin-gonic/gin"
)

func Hello(ctx *gin.Context) {
	util.Success(ctx, "Hello, World!")
}
