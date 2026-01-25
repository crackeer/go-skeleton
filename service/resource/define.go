package resource

import "io"

type Entry struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Type       string `json:"type"`
	ModifyTime int64  `json:"modify_time"`
}

type Resource interface {
	List(string) ([]Entry, error)
	Read(string) (io.Reader, error)
	Write(string, io.Reader) error
	Delete(string) error
	MkdirAll(string) error
}
