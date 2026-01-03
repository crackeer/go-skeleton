package util

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type JSONResponse struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// Success  ...
//
//	@param ctx
//	@param data
func Success(ctx *gin.Context, data interface{}) {
	ctx.AbortWithStatusJSON(http.StatusOK, &JSONResponse{
		Code:    0,
		Message: "success",
		Data:    data,
	})
}

// Failure ...
//
//	@param ctx
//	@param code
//	@param message
func Failure(ctx *gin.Context, code int64, message string) {
	ctx.AbortWithStatusJSON(http.StatusOK, &JSONResponse{
		Code:    code,
		Message: message,
		Data:    nil,
	})
}
