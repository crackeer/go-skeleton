package template

import (
	_ "embed"

	"github.com/crackeer/go-connect/container"
	"github.com/crackeer/go-connect/service/resource"
	"github.com/flosch/pongo2/v7"
)

var (
	//go:embed list.template
	listTemplate string
)

// RenderList 使用pongo2渲染文件列表模板
func RenderList(entries []resource.Entry, title string, path string, currentURL string) string {
	var (
		tpl *pongo2.Template
		err error
	)
	if container.IsDevelop() {
		tpl, err = pongo2.FromFile("./service/template/list.template")
	} else {
		tpl, err = pongo2.FromString(listTemplate)
	}

	if err != nil {
		return "Failed to load template: " + err.Error()
	}

	// 准备上下文数据
	ctx := pongo2.Context{
		"entries": entries,
		"title":   title,
		"path":    path,
		"url":     currentURL,
	}

	// 渲染模板
	result, err := tpl.Execute(ctx)
	if err != nil {
		return "Failed to render template: " + err.Error()
	}

	return result
}
