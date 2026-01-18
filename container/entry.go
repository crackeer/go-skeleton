package container

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	Name          string `yaml:"name"`
	Driver        string `yaml:"driver"`
	File          string `yaml:"file,omitempty"`
	WriteHost     string `yaml:"write_host,omitempty"`
	ReadHost      string `yaml:"read_host,omitempty"`
	WriteUser     string `yaml:"write_user,omitempty"`
	ReadUser      string `yaml:"read_user,omitempty"`
	WritePassword string `yaml:"write_password,omitempty"`
	ReadPassword  string `yaml:"read_password,omitempty"`
	Charset       string `yaml:"charset,omitempty"`
	Database      string `yaml:"database,omitempty"`
}

// AppConfig Config
type AppConfig struct {
	Port            int64            `yaml:"port"`
	DefaultPageSize int64            `yaml:"default_page_size"`
	PublicDir       string           `yaml:"public_dir"`
	Env             string           `yaml:"env"`
	Database        []DatabaseConfig `yaml:"database"`
	Redis           []RedisConfig    `yaml:"redis"`
}

type RedisConfig struct {
	Name     string `yaml:"name"`
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	DB       int    `yaml:"db"`
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
	if len(appConfig.Database) > 0 {
		if err := initDatabase(appConfig.Database); err != nil {
			return nil, fmt.Errorf("init database connection error: %s", err.Error())
		}
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
