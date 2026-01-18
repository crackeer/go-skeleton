package api

import (
	"fmt"
	"strings"

	"github.com/crackeer/go-connect/service/template"
	"github.com/crackeer/go-connect/util"

	"github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	config, err := GetDriverConfig(name)
	if err != nil {
		util.Failure(ctx, -1, err.Error())
		return
	}
	path := strings.Trim(ctx.Param("path"), "/")
	client, err := NewResourceClient(config)
	list, err := client.List(path)
	if err != nil {
		util.Failure(ctx, -1, err.Error())
		return
	}
	fmt.Println(list)

	// 渲染模板
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(200, template.RenderList(list, config.Title, path, strings.Trim(ctx.Request.URL.Path, "/")))
}
