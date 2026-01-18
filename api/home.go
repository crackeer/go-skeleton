package api

import (
	"net/http"

	"github.com/crackeer/go-connect/container"
	"github.com/crackeer/go-connect/service/template"

	"github.com/gin-gonic/gin"
)

func Home(ctx *gin.Context) {
	// 获取所有的资源配置
	appConfig := container.GetAppConfig()

	// 渲染资源列表模板
	html := template.RenderHome(appConfig.Resource)

	// 返回HTML响应
	ctx.Header("Content-Type", "text/html; charset=utf-8")
	ctx.String(http.StatusOK, html)
}
