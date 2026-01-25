package resource

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)

type LocalDriver struct {
	RootDir string
}

func NewLocalDriver(rootDir string) (*LocalDriver, error) {
	fmt.Println(rootDir)
	return &LocalDriver{RootDir: rootDir}, nil
}

// List 列出指定路径下的文件和子目录
func (d *LocalDriver) List(path string) ([]Entry, error) {
	fullPath := filepath.Join(d.RootDir, path)
	fmt.Println(fullPath)
	entries, err := os.ReadDir(fullPath)

	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	var result []Entry
	for _, entry := range entries {
		info, err := entry.Info()
		if err != nil {
			continue
		}

		entryType := "file"
		if entry.IsDir() {
			entryType = "dir"
		}

		result = append(result, Entry{
			Name:       entry.Name(),
			Size:       info.Size(),
			Type:       entryType,
			ModifyTime: info.ModTime().Unix(),
		})
	}

	return result, nil
}

// Read 读取指定文件的内容
func (d *LocalDriver) Read(path string) (io.Reader, error) {
	fullPath := filepath.Join(d.RootDir, path)
	return os.Open(fullPath)
}

// Write 将数据写入指定文件
func (d *LocalDriver) Write(path string, data io.Reader) error {
	fullPath := filepath.Join(d.RootDir, path)

	// 确保目录结构存在
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}

	// 写入文件
	out, err := os.Create(fullPath)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, data); err != nil {
		return err
	}

	return nil
}

// Delete 删除指定路径的文件或目录
func (d *LocalDriver) Delete(path string) error {
	fullPath := filepath.Join(d.RootDir, path)

	// 检查路径是文件还是目录
	info, err := os.Stat(fullPath)
	if err != nil {
		return err
	}

	if info.IsDir() {
		// 删除目录
		return os.RemoveAll(fullPath)
	} else {
		// 删除文件
		return os.Remove(fullPath)
	}
}

// MkdirAll 递归创建目录结构
func (d *LocalDriver) MkdirAll(path string) error {
	fullPath := filepath.Join(d.RootDir, path)
	return os.MkdirAll(fullPath, 0755)
}
