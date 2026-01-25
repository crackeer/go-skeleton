package resource

import (
	"bytes"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/jlaffaye/ftp"
)

type FTPDriver struct {
	Host        string
	Port        int
	Username    string
	Password    string
	RelativeDir string
}

func NewFTPDriver(host string, port int, username string, password string, relativeDir string) (*FTPDriver, error) {
	return &FTPDriver{
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		RelativeDir: relativeDir,
	}, nil
}

// connect 建立FTP连接
func (d *FTPDriver) connect() (*ftp.ServerConn, error) {
	// 连接FTP服务器
	addr := fmt.Sprintf("%s:%d", d.Host, d.Port)
	conn, err := ftp.Dial(addr, ftp.DialWithTimeout(5*time.Second))
	if err != nil {
		return nil, err
	}

	// 登录FTP服务器
	if err := conn.Login(d.Username, d.Password); err != nil {
		conn.Quit()
		return nil, err
	}

	// 设置相对目录
	if d.RelativeDir != "" {
		if err := conn.ChangeDir(d.RelativeDir); err != nil {
			conn.Quit()
			return nil, err
		}
	}

	return conn, nil
}

// List 列出指定路径下的文件和子目录
func (d *FTPDriver) List(path string) ([]Entry, error) {
	fmt.Println(path, d.RelativeDir)
	// 建立FTP连接
	conn, err := d.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Quit()
	fmt.Println(path)
	// 列出文件和目录
	entries, err := conn.List(path)
	if err != nil {
		return nil, err
	}

	var result []Entry
	for _, entry := range entries {
		// 跳过.和..
		if entry.Name == "." || entry.Name == ".." {
			continue
		}

		entryType := "file"
		if entry.Type == ftp.EntryTypeFolder {
			entryType = "dir"
		}

		result = append(result, Entry{
			Name:       entry.Name,
			Size:       int64(entry.Size),
			Type:       entryType,
			ModifyTime: entry.Time.Unix(),
		})
	}

	return result, nil
}

// Read 读取指定文件的内容
func (d *FTPDriver) Read(path string) (io.Reader, error) {
	// 建立FTP连接
	conn, err := d.connect()
	if err != nil {
		return nil, err
	}
	defer conn.Quit()

	// 获取文件内容
	data, err := conn.Retr(path)
	if err != nil {
		return nil, err
	}
	defer data.Close()

	// 读取数据
	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(data)
	if err != nil {
		return nil, err
	}

	return buf, nil
}

// makeDirAll 递归创建目录结构
func (d *FTPDriver) makeDirAll(conn *ftp.ServerConn, dir string) error {
	// 将路径分割成目录组件
	dirs := strings.Split(dir, string(filepath.Separator))
	currentPath := ""

	// 递归创建每个目录
	for _, d := range dirs {
		if d == "" {
			continue
		}

		currentPath = filepath.Join(currentPath, d)
		if err := conn.ChangeDir(currentPath); err != nil {
			// 如果目录不存在，创建它
			if err := conn.MakeDir(currentPath); err != nil {
				return err
			}
			// 切换到新创建的目录
			if err := conn.ChangeDir(currentPath); err != nil {
				return err
			}
		}
	}

	// 切换回根目录
	return conn.ChangeDir("/")
}

// Write 将数据写入指定文件
func (d *FTPDriver) Write(path string, data io.Reader) error {
	// 建立FTP连接
	conn, err := d.connect()
	if err != nil {
		return err
	}
	defer conn.Quit()

	// 确保目录结构存在
	dir := filepath.Dir(path)
	if dir != "." && dir != "" {
		if err := d.makeDirAll(conn, dir); err != nil {
			return err
		}
	}

	// 写入文件
	return conn.Stor(path, data)
}

// Delete 删除指定路径的文件或目录
func (d *FTPDriver) Delete(path string) error {
	// 建立FTP连接
	conn, err := d.connect()
	if err != nil {
		return err
	}
	defer conn.Quit()

	// 尝试删除目录
	if err := conn.RemoveDir(path); err == nil {
		// 目录删除成功
		return nil
	}

	// 尝试删除文件
	if err := conn.Delete(path); err == nil {
		// 文件删除成功
		return nil
	}

	// 无法确定是文件还是目录，或者删除失败
	return fmt.Errorf("failed to delete %s", path)
}

// MkdirAll 递归创建目录结构
func (d *FTPDriver) MkdirAll(path string) error {
	// 建立FTP连接
	conn, err := d.connect()
	if err != nil {
		return err
	}
	defer conn.Quit()

	// 使用现有的makeDirAll方法创建目录结构
	return d.makeDirAll(conn, path)
}
