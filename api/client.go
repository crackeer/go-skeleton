package api

import (
	"errors"

	"github.com/crackeer/go-connect/container"
	"github.com/crackeer/go-connect/service/resource"
)

func GetDriverConfig(name string) (container.DriverConfig, error) {
	appConfig := container.GetAppConfig()
	var (
		index int = -1
	)
	for i, item := range appConfig.Resource {
		if item.Name == name {
			index = i
			break
		}
	}
	if index < 0 {
		return container.DriverConfig{}, errors.New("resource not found")
	}
	return appConfig.Resource[index], nil
}

func NewResourceClient(config container.DriverConfig) (resource.Resource, error) {
	var (
		rsc resource.Resource
		err error
	)
	switch config.Driver {
	case "local":
		rsc, err = resource.NewLocalDriver(config.LocalConfig.Dir)
	case "s3":
		rsc, err = resource.NewS3Driver(config.S3Config.AccessKey, config.S3Config.SecretKey, config.S3Config.Bucket, config.S3Config.Endpoint, config.S3Config.Region)
	case "ftp":
		rsc, err = resource.NewFTPDriver(config.FTPConfig.Host, int(config.FTPConfig.Port), config.FTPConfig.User, config.FTPConfig.Password, config.FTPConfig.RelativePath)
	default:
		return nil, errors.New("resource driver not found")
	}
	if err != nil {
		return nil, err
	}

	return rsc, nil
}
