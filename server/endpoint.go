package server

import (
	"fmt"

	"github.com/crackeer/go-connect/api"
	"github.com/crackeer/go-connect/container"

	"github.com/gin-gonic/gin"
)

// Run ...
//
//	@param config
//	@return error
func Run(config *container.AppConfig) error {
	router := gin.New()
	router.Use(gin.Logger())
	if len(config.User) > 0 {
		router.Use(gin.BasicAuth(gin.Accounts(toGinAccounts(config.User))))
	}
	router.GET("/:name/*path", api.Get)
	//router.POST("/upload/:name/*path", api.Upload)
	router.POST("/:name/*path", api.Upload)
	router.DELETE("/:name/*path", api.Delete)
	router.NoRoute(api.Home)
	return router.Run(fmt.Sprintf(":%d", config.Port))
}

func toGinAccounts(users []container.User) gin.Accounts {
	accounts := make(gin.Accounts)
	for _, user := range users {
		accounts[user.Name] = user.Password
	}
	return accounts
}
