package server

import (
	"fmt"
	"go-skeleton/api"

	"github.com/gin-gonic/gin"
)

// Run ...
//
//	@param config
//	@return error
func Run(port int64) error {
	router := gin.New()
	router.GET("/hello", api.Hello)
	return router.Run(fmt.Sprintf(":%d", port))
}
