package adaptors

import (
	"io"
	"os"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/pkg/errors"
)

func NewFileSystem() adaptors.FileSystem {
	return &FileSystem{}
}

type FileSystem struct{}

func (*FileSystem) Open(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return f, nil
}

func (*FileSystem) Create(name string) (io.WriteCloser, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return f, nil
}

func (*FileSystem) Remove(name string) error {
	if err := os.Remove(name); err != nil {
		return errors.WithStack(err)
	}
	return nil
}

func (*FileSystem) Stat(name string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return fileInfo, nil
}

func (*FileSystem) MkdirAll(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
