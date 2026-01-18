package api

import (
	"strings"

	"github.com/crackeer/go-connect/util"

	"github.com/gin-gonic/gin"
)

func Upload(ctx *gin.Context) {
	// 获取驱动名称
	driverName := ctx.Param("name")

	// 获取上传的文件
	file, err := ctx.FormFile("file")
	if err != nil {
		util.Failure(ctx, -1, "获取文件失败: "+err.Error())
		return
	}

	// 打开文件
	src, err := file.Open()
	if err != nil {
		util.Failure(ctx, -1, "打开文件失败: "+err.Error())
		return
	}
	defer src.Close()

	// 创建文件内容缓冲区
	buffer := make([]byte, file.Size)
	_, err = src.Read(buffer)
	if err != nil {
		util.Failure(ctx, -1, "读取文件失败: "+err.Error())
		return
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

	// 从URL获取路径
	path := ctx.Param("path")

	// 写入文件
	err = client.Write(strings.Trim(path, "/"), buffer)
	if err != nil {
		util.Failure(ctx, -1, "写入文件失败: "+err.Error())
		return
	}

	util.Success(ctx, gin.H{"message": "文件上传成功"})
}
