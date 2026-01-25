# go-connect

一个基于 Go 语言开发的文件管理系统，支持本地文件、S3 对象存储和 FTP 服务器的统一管理。

## 功能特性

- **多存储驱动支持**：
  - 本地文件系统
  - S3 对象存储（支持兼容 S3 协议的存储服务）
  - FTP 服务器
- **文件管理功能**：
  - 浏览文件和目录
  - 上传文件（支持批量上传）
  - 下载文件
  - 删除文件（支持批量删除）
  - 新建文件夹
- **用户认证**：支持基本认证
- **Web 界面**：友好的 Web 操作界面

## 技术栈

- **后端**：Go 1.25
- **Web 框架**：Gin
- **模板引擎**：Pongo2
- **存储驱动**：
  - 本地文件系统：Go 标准库
  - S3 对象存储：AWS SDK for Go v2
  - FTP 服务器：jlaffaye/ftp
- **配置管理**：YAML

## 项目结构

```
go-connect/
├── api/              # API 端点实现
│   ├── client.go     # 资源客户端创建
│   ├── delete.go     # 删除文件 API
│   ├── get.go        # 获取文件/目录列表 API
│   ├── home.go       # 首页 API
│   ├── mkdir.go      # 创建目录 API
│   └── upload.go     # 上传文件 API
├── config/           # 配置文件
│   └── app.yaml      # 应用配置
├── container/        # 依赖注入容器
│   └── entry.go      # 容器入口
├── server/           # 服务器配置
│   └── endpoint.go   # 路由配置
├── service/          # 业务逻辑
│   ├── resource/     # 存储驱动实现
│   │   ├── define.go # 资源接口定义
│   │   ├── ftp.go    # FTP 驱动
│   │   ├── local.go  # 本地文件系统驱动
│   │   └── s3.go     # S3 驱动
│   └── template/     # 模板渲染
│       ├── entry.go  # 模板渲染入口
│       ├── home.html # 首页模板
│       └── list.html # 文件列表模板
├── util/             # 工具函数
│   └── gin.go        # Gin 相关工具
├── go.mod            # Go 模块文件
├── go.sum            # Go 依赖校验文件
├── main.go           # 应用入口
└── README.md         # 项目说明
```

## 快速开始

### 前置条件

- Go 1.25 或更高版本
- （可选）S3 兼容的存储服务
- （可选）FTP 服务器

### 安装

1. 克隆项目代码：

```bash
git clone https://github.com/crackeer/go-skeleton.git
cd go-skeleton
```

2. 安装依赖：

```bash
go mod download
```

### 配置

编辑 `config/app.yaml` 文件，配置存储服务和认证信息：

```yaml
port: 9000
default_page_size: 20
env: "develop"
user:
  - name: "admin"
    password: "admin"
resource:
  - driver: "local"
    name: "local"
    title: "本地文件系统"
    local_config:
      dir: "./"
  - driver: "s3"
    name: "s3"
    title: "S3对象存储"
    s3_config:
      bucket: "test"
      region: "us-east-1"
      access_key: "test"
      secret_key: "test"
      endpoint: "http://localhost:9000"
  - driver: "ftp"
    name: "ftp"
    title: "FTP文件系统"
    ftp_config:
      host: "10.33.202.152"
      port: 8090
      user: "ftpuser"
      password: "ftpuser"
      relative_path: "/"
```

### 启动

1. 编译项目：

```bash
go build -o go-connect .
```

2. 运行项目：

```bash
./go-connect
```

或者直接运行：

```bash
go run main.go
```

3. 访问 Web 界面：

打开浏览器，访问 `http://localhost:9000`，使用配置文件中的用户名和密码登录。

## API 文档

### 基础路径

所有 API 端点都以 `/<驱动名称>/<路径>` 格式开始，其中 `<驱动名称>` 是配置文件中定义的资源名称。

### 端点

- **GET /<驱动名称>/<路径>**：获取文件或目录列表
  - 参数：
    - `download=true`：下载文件（仅对文件有效）

- **PUT /<驱动名称>/<路径>**：上传文件
  - 表单数据：
    - `file`：要上传的文件

- **POST /<驱动名称>/<路径>**：创建目录

- **DELETE /<驱动名称>/<路径>**：删除文件或目录

## 示例

### 1. 浏览本地文件系统

访问 `http://localhost:9000/local` 查看本地文件系统的根目录。

### 2. 上传文件到 S3

1. 访问 `http://localhost:9000/s3`
2. 点击 "上传文件" 按钮
3. 选择要上传的文件
4. 等待上传完成

### 3. 创建目录

1. 访问任意存储驱动的目录
2. 点击 "新建文件夹" 按钮
3. 输入文件夹名称
4. 点击 "确定" 按钮

### 4. 批量删除文件

1. 访问任意存储驱动的目录
2. 勾选要删除的文件或目录
3. 点击 "批量删除" 按钮
4. 确认删除操作

## 配置说明

### 服务器配置

- `port`：服务器监听端口
- `env`：环境变量（`develop` 或 `production`）
- `default_page_size`：默认分页大小

### 用户配置

- `user`：用户列表
  - `name`：用户名
  - `password`：密码

### 资源配置

#### 本地文件系统

- `driver`：`local`
- `name`：资源名称（唯一）
- `title`：显示名称
- `local_config.dir`：本地目录路径

#### S3 对象存储

- `driver`：`s3`
- `name`：资源名称（唯一）
- `title`：显示名称
- `s3_config.bucket`：存储桶名称
- `s3_config.region`：区域
- `s3_config.access_key`：访问密钥
- `s3_config.secret_key`：秘密密钥
- `s3_config.endpoint`：S3 服务端点（可选，默认为 AWS S3 端点）

#### FTP 服务器

- `driver`：`ftp`
- `name`：资源名称（唯一）
- `title`：显示名称
- `ftp_config.host`：FTP 服务器主机
- `ftp_config.port`：FTP 服务器端口
- `ftp_config.user`：FTP 用户名
- `ftp_config.password`：FTP 密码
- `ftp_config.relative_path`：相对路径

## 开发

### 代码风格

遵循 Go 语言标准代码风格，使用 `gofmt` 格式化代码。

### 构建和测试

```bash
# 构建
make build

# 运行
make run
```

## 许可证

本项目采用 MIT 许可证。

## 贡献

欢迎提交 Issue 和 Pull Request 来改进这个项目。

## 联系方式

- 项目地址：https://github.com/crackeer/go-skeleton
- 作者：crackeer