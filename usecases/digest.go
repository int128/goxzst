package usecases

import (
	"crypto/sha256"
	"fmt"
	"io"
	"path/filepath"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewDigest(i Digest) usecases.Digest {
	return &i
}

type Digest struct {
	dig.In
	Filesystem adaptors.Filesystem
	Logger     adaptors.Logger
}

func (u *Digest) Do(in usecases.DigestIn) (*usecases.DigestOut, error) {
	if err := u.Filesystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return nil, errors.Wrapf(err, "error while creating the output directory")
	}

	u.Logger.Logf("Creating %s", in.OutputFilename)
	input, err := u.Filesystem.Open(in.InputFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "error while opening the file %s", in.InputFilename)
	}
	defer input.Close()

	w := sha256.New()
	if _, err := io.Copy(w, input); err != nil {
		return nil, errors.Wrapf(err, "error while computing digest of the file %s", in.InputFilename)
	}
	h := fmt.Sprintf("%x", w.Sum(nil))

	output, err := u.Filesystem.Create(in.OutputFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating the file %s", in.OutputFilename)
	}
	defer output.Close()
	n, err := output.Write([]byte(h))
	if err != nil {
		return nil, errors.Wrapf(err, "error while writing to the file %s", in.OutputFilename)
	}
	if n != len(h) {
		return nil, errors.Errorf("wants to write %d bytes but wrote %d bytes to the file %s", len(h), n, in.OutputFilename)
	}

	return &usecases.DigestOut{
		SHA256: h,
	}, nil
}
