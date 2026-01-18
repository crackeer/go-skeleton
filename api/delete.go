package api

import (
	"github.com/crackeer/go-connect/util"

	"github.com/gin-gonic/gin"
)

func Delete(ctx *gin.Context) {
	// 获取驱动名称
	driverName := ctx.Param("name")

	// 从URL获取路径
	path := ctx.Param("path")
	// 去除路径开头的斜杠
	if len(path) > 0 && path[0] == '/' {
		path = path[1:]
	}

	// 获取驱动配置
	driverConfig, err := GetDriverConfig(driverName)
	if err != nil {
		util.Failure(ctx, -1, "获取驱动配置失败: "+err.Error())
		return
	}

	// 创建资源客户端
	client, err := NewResourceClient(driverConfig)
	if err != nil {
		util.Failure(ctx, -1, "创建客户端失败: "+err.Error())
		return
	}

	// 删除文件或目录
	err = client.Delete(path)
	if err != nil {
		util.Failure(ctx, -1, "删除失败: "+err.Error())
		return
	}

	util.Success(ctx, gin.H{"message": "删除成功"})
}
