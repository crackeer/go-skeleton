package container

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

const (
	DriverLocal = "local"

	DriverS3 = "s3"

	DriverFTP = "ftp"
)

const (
	EntryTypeFile = "file"

	EntryTypeFolder = "folder"
)

type LocalConfig struct {
	Dir string
}

type S3Config struct {
	Bucket    string `yaml:"bucket"`
	Region    string `yaml:"region"`
	AccessKey string `yaml:"access_key"`
	SecretKey string `yaml:"secret_key"`
	Endpoint  string `yaml:"endpoint"`
}

type FTPConfig struct {
	Host         string `yaml:"host"`
	Port         int64  `yaml:"port"`
	User         string `yaml:"user"`
	Password     string `yaml:"password"`
	RelativePath string `yaml:"relative_path"`
}

type DriverConfig struct {
	Driver      string      `yaml:"driver"`
	Name        string      `yaml:"name"`
	LocalConfig LocalConfig `yaml:"local_config"`
	S3Config    S3Config    `yaml:"s3_config"`
	FTPConfig   FTPConfig   `yaml:"ftp_config"`
	Title       string      `yaml:"title"`
}

type User struct {
	Name     string `json:"name"`
	Password string `json:"password"`
}

// AppConfig Config
type AppConfig struct {
	Port            int64          `yaml:"port"`
	Env             string         `yaml:"env"`
	DefaultPageSize int64          `yaml:"default_page_size"`
	User            []User         `yaml:"user"`
	Resource        []DriverConfig `yaml:"resource"`
}

var appConfig *AppConfig

// InitConfig ...
//
//	@param configPath
//	@return error
func Init(configPath string) (*AppConfig, error) {
	appConfig = &AppConfig{}
	if err := loadYaml(configPath, appConfig); err != nil {
		return nil, fmt.Errorf("load app config error: %s", err.Error())
	}
	return appConfig, nil
}

// GetAppConfig ...
//
//	@return *define.AppConfig
func GetAppConfig() *AppConfig {
	return appConfig
}

func loadYaml(path string, dest interface{}) error {
	bytes, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read config `%s` error: %s", path, err.Error())
	}

	return yaml.Unmarshal(bytes, dest)
}

// IsDevelop ...
//
//	@return bool
func IsDevelop() bool {
	return appConfig.Env == "develop"
}
