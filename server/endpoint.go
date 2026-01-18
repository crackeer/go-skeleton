package server

import (
	"fmt"
	"net/http"

	"github.com/crackeer/go-skeleton/api"
	"github.com/crackeer/go-skeleton/container"

	"github.com/gin-gonic/gin"
)

// Run ...
//
//	@param config
//	@return error
func Run(config *container.AppConfig) error {
	router := gin.New()
	router.GET("/hello", api.Hello)
	router.StaticFS("/public", http.Dir(config.PublicDir))
	return router.Run(fmt.Sprintf(":%d", config.Port))
}
