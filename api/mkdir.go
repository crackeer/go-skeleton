package api

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func Mkdir(ctx *gin.Context) {
	// 获取驱动名称
	driverName := ctx.Param("name")

	// 获取驱动配置
	driverConfig, err := GetDriverConfig(driverName)
	if err != nil {
		ctx.String(http.StatusBadRequest, "获取驱动配置失败: "+err.Error())
		return
	}

	// 创建资源客户端
	client, err := NewResourceClient(driverConfig)
	if err != nil {
		ctx.String(http.StatusBadRequest, "创建客户端失败: "+err.Error())
		return
	}

	// 从URL获取路径
	path := ctx.Param("path")

	// 创建目录
	err = client.MkdirAll(strings.Trim(path, "/"))
	if err != nil {
		ctx.String(http.StatusBadRequest, "创建目录失败: "+err.Error())
		return
	}

	ctx.String(http.StatusOK, "目录创建成功")
}
