package adaptors

import (
	"io"
	"os"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/pkg/errors"
)

func NewFilesystem() adaptors.Filesystem {
	return &Filesystem{}
}

type Filesystem struct{}

func (*Filesystem) Open(name string) (io.ReadCloser, error) {
	f, err := os.Open(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return f, nil
}

func (*Filesystem) Create(name string) (io.WriteCloser, error) {
	f, err := os.Create(name)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return f, nil
}

func (*Filesystem) GetMode(name string) (os.FileMode, error) {
	s, err := os.Stat(name)
	if err != nil {
		return 0, errors.WithStack(err)
	}
	return s.Mode(), nil
}

func (*Filesystem) MkdirAll(path string) error {
	if err := os.MkdirAll(path, 0755); err != nil {
		return errors.WithStack(err)
	}
	return nil
}
