// Package digest provides the use-case to generate a digest.
package digest

import (
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/models/digest"
)

var Set = wire.NewSet(
	wire.Struct(new(Digest), "*"),
	wire.Bind(new(Interface), new(*Digest)),
)

//go:generate mockgen -destination mock_digest/mock_digest.go github.com/int128/goxzst/usecases/digest Interface

type Interface interface {
	Do(in Input) error
}

type Input struct {
	InputFilename  string
	OutputFilename string
	Algorithm      *digest.Algorithm
}

type Digest struct {
	FileSystem fs.Interface
	Logger     logger.Interface
}

func (u *Digest) Do(in Input) error {
	if err := u.FileSystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return fmt.Errorf("error while creating the output directory: %w", err)
	}

	u.Logger.Logf("Creating %s", in.OutputFilename)
	input, err := u.FileSystem.Open(in.InputFilename)
	if err != nil {
		return fmt.Errorf("error while opening the file %s: %w", in.InputFilename, err)
	}
	defer input.Close()

	w := in.Algorithm.NewHash()
	if _, err := io.Copy(w, input); err != nil {
		return fmt.Errorf("error while computing digest of the file %s: %w", in.InputFilename, err)
	}
	h := fmt.Sprintf("%x", w.Sum(nil))

	output, err := u.FileSystem.Create(in.OutputFilename)
	if err != nil {
		return fmt.Errorf("error while creating the file %s: %w", in.OutputFilename, err)
	}
	defer output.Close()
	n, err := output.Write([]byte(h))
	if err != nil {
		return fmt.Errorf("error while writing to the file %s: %w", in.OutputFilename, err)
	}
	if n != len(h) {
		return fmt.Errorf("wants to write %d bytes but wrote %d bytes to the file %s", len(h), n, in.OutputFilename)
	}
	return nil
}
