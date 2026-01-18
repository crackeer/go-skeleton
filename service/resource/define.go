package resource

type Entry struct {
	Name       string `json:"name"`
	Size       int64  `json:"size"`
	Type       string `json:"type"`
	ModifyTime int64  `json:"modify_time"`
}

type Resource interface {
	List(string) ([]Entry, error)
	Detail(string) (Entry, error)
	Read(string) ([]byte, error)
	Write(string, []byte) error
	Delete(string) error
}
