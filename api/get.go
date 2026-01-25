package api

import (
	"net/http"
	"path/filepath"
	"strings"

	"github.com/crackeer/go-connect/service/template"

	"github.com/gin-gonic/gin"
)

func Get(ctx *gin.Context) {
	name := ctx.Param("name")
	config, err := GetDriverConfig(name)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	path := strings.Trim(ctx.Param("path"), "/")
	client, err := NewResourceClient(config)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}
	if ctx.Query("download") == "true" {
		data, err := client.Read(path)
		if err != nil {
			ctx.String(http.StatusOK, err.Error())
			return
		}
		ctx.DataFromReader(http.StatusOK, 0, filepath.Base(path), data, map[string]string{
			"Content-Disposition": "attachment; filename=" + filepath.Base(path),
		})
		return
	}
	list, err := client.List(path)
	if err != nil {
		ctx.String(http.StatusOK, err.Error())
		return
	}

	// 渲染模板
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(200, template.RenderList(list, config.Title, path, strings.Trim(ctx.Request.URL.Path, "/")))
}
