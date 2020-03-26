package fs

import (
	"fmt"
	"io"
	"os"

	"github.com/google/wire"
)

var Set = wire.NewSet(
	wire.Struct(new(FileSystem), "*"),
	wire.Bind(new(Interface), new(*FileSystem)),
)

//go:generate mockgen -destination mock_fs/mock_fs.go github.com/int128/goxzst/adaptors/fs Interface

type Interface interface {
	Open(name string) (io.ReadCloser, error)
	Create(name string) (io.WriteCloser, error)
	Remove(name string) error
	Stat(name string) (os.FileInfo, error)
	MkdirAll(path string) error
}

type FileSystem struct{}

func (*FileSystem) Open(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return f, nil
}

func (*FileSystem) Create(name string) (io.WriteCloser, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return f, nil
}

func (*FileSystem) Remove(name string) error {
	if err := os.Remove(name); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}

func (*FileSystem) Stat(name string) (os.FileInfo, error) {
	fileInfo, err := os.Stat(name)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}
	return fileInfo, nil
}

func (*FileSystem) MkdirAll(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("%w", err)
	}
	return nil
}
