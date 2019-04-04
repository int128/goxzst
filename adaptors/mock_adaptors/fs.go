package mock_adaptors

import (
	"os"
	"time"
)

type FileInfo struct {
	ModeValue    os.FileMode
	ModTimeValue time.Time
}

func (f *FileInfo) Name() string {
	panic("Name() not implemented")
}

func (f *FileInfo) Size() int64 {
	panic("Size() not implemented")
}

func (f *FileInfo) Mode() os.FileMode {
	return f.ModeValue
}

func (f *FileInfo) ModTime() time.Time {
	return f.ModTimeValue
}

func (f *FileInfo) IsDir() bool {
	panic("IsDir() not implemented")
}

func (f *FileInfo) Sys() interface{} {
	panic("Sys() not implemented")
}
