#!/bin/bash

default_mod="github.com/crackeer/go-skeleton"

# 检查是否提供了新的mod参数
if [ $# -eq 0 ]; then
    echo "Usage: $0 <new_mod_path>"
    exit 1
fi

new_mod=$1

# 替换go.mod中的module声明
sed -i "s/module $default_mod/module $new_mod/g" go.mod

# 更新所有.go文件中的导入路径
grep -r "$default_mod/" --include="*.go" . | cut -d: -f1 | sort -u | xargs sed -i "s,$default_mod/,$new_mod/,g"

# 运行go mod tidy
echo "Running go mod tidy..."
go mod tidy

echo "Module updated from $default_mod to $new_mod successfully!"
