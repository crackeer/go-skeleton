package resource

import (
	"bytes"
	"context"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Driver struct {
	AccessKey string
	SecretKey string
	Bucket    string
	Endpoint  string
	Region    string
	Client    *s3.Client
}

func NewS3Driver(accessKey, secretKey, bucket, endpoint, region string) (*S3Driver, error) {
	// 创建凭证提供者
	creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")

	// 加载配置
	cfg, err := config.LoadDefaultConfig(context.Background(),
		config.WithCredentialsProvider(creds),
		config.WithRegion(region),
	)
	if err != nil {
		return nil, err
	}

	// 创建S3客户端
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		if endpoint != "" {
			o.BaseEndpoint = aws.String(endpoint)
		}
	})

	return &S3Driver{
		AccessKey: accessKey,
		SecretKey: secretKey,
		Bucket:    bucket,
		Endpoint:  endpoint,
		Region:    region,
		Client:    client,
	}, nil
}

// List 列出指定路径下的对象
func (d *S3Driver) List(path string) ([]Entry, error) {
	// 确保路径以斜杠结尾（如果是目录）
	if path != "" && !strings.HasSuffix(path, "/") {
		path += "/"
	}

	// 列出对象
	resp, err := d.Client.ListObjectsV2(context.Background(), &s3.ListObjectsV2Input{
		Bucket:    aws.String(d.Bucket),
		Prefix:    aws.String(path),
		Delimiter: aws.String("/"),
	})
	if err != nil {
		return nil, err
	}

	var result []Entry

	// 处理目录（CommonPrefixes）
	for _, prefix := range resp.CommonPrefixes {
		dirName := filepath.Base(*prefix.Prefix)
		result = append(result, Entry{
			Name:       dirName,
			Size:       0,
			Type:       "dir",
			ModifyTime: time.Now().Unix(),
		})
	}

	// 处理文件
	for _, obj := range resp.Contents {
		if *obj.Key == path { // 跳过当前路径自身
			continue
		}

		fileName := filepath.Base(*obj.Key)
		result = append(result, Entry{
			Name:       fileName,
			Size:       *obj.Size,
			Type:       "file",
			ModifyTime: obj.LastModified.Unix(),
		})
	}

	return result, nil
}

// Read 读取指定对象的内容
func (d *S3Driver) Read(path string) ([]byte, error) {
	// 获取对象
	resp, err := d.Client.GetObject(context.Background(), &s3.GetObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// 读取对象内容
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}

// Write 将数据写入指定对象
func (d *S3Driver) Write(path string, data []byte) error {
	// 写入对象
	_, err := d.Client.PutObject(context.Background(), &s3.PutObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return err
	}

	return nil
}

// Detail 获取指定对象的详细信息
func (d *S3Driver) Detail(path string) (Entry, error) {
	// 获取对象元数据
	resp, err := d.Client.HeadObject(context.Background(), &s3.HeadObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return Entry{}, err
	}

	// 获取对象名称
	name := filepath.Base(path)

	// 确定对象类型
	entryType := "file"
	// 如果路径以斜杠结尾，或者大小为0且没有Content-Type，可能是目录
	isZeroSize := resp.ContentLength != nil && *resp.ContentLength == 0
	if strings.HasSuffix(path, "/") || (isZeroSize && resp.ContentType == nil) {
		entryType = "dir"
	}

	return Entry{
		Name:       name,
		Size:       *resp.ContentLength,
		Type:       entryType,
		ModifyTime: resp.LastModified.Unix(),
	}, nil
}

// Delete 删除指定路径的对象
func (d *S3Driver) Delete(path string) error {
	// 如果是文件，直接删除
	_, err := d.Client.DeleteObject(context.Background(), &s3.DeleteObjectInput{
		Bucket: aws.String(d.Bucket),
		Key:    aws.String(path),
	})
	return err
}
