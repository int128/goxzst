package digest

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/int128/goxzst/adaptors"
	"github.com/int128/goxzst/usecases"
	"github.com/pkg/errors"
)

type Digest struct {
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
}

func (u *Digest) Do(in usecases.DigestIn) error {
	if err := u.FileSystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return errors.Wrapf(err, "error while creating the output directory")
	}

	u.Logger.Logf("Creating %s", in.OutputFilename)
	input, err := u.FileSystem.Open(in.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while opening the file %s", in.InputFilename)
	}
	defer input.Close()

	w := in.Algorithm.NewHash()
	if _, err := io.Copy(w, input); err != nil {
		return errors.Wrapf(err, "error while computing digest of the file %s", in.InputFilename)
	}
	h := fmt.Sprintf("%x", w.Sum(nil))

	output, err := u.FileSystem.Create(in.OutputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while creating the file %s", in.OutputFilename)
	}
	defer output.Close()
	n, err := output.Write([]byte(h))
	if err != nil {
		return errors.Wrapf(err, "error while writing to the file %s", in.OutputFilename)
	}
	if n != len(h) {
		return errors.Errorf("wants to write %d bytes but wrote %d bytes to the file %s", len(h), n, in.OutputFilename)
	}
	return nil
}
