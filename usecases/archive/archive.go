// Package archive provides the use-case to archive an executable file and extra files.
package archive

import (
	"archive/zip"
	"fmt"
	"io"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
)

var Set = wire.NewSet(
	wire.Struct(new(Archive), "*"),
	wire.Bind(new(Interface), new(*Archive)),
)

//go:generate mockgen -destination mock_archive/mock_archive.go github.com/int128/goxzst/usecases/archive Interface

type Interface interface {
	Do(in Input) error
}

type Input struct {
	OutputFilename string
	Entries        []Entry
}

type Entry struct {
	Filename      string // filename in the archive
	InputFilename string
}

type Archive struct {
	FileSystem fs.Interface
	Logger     logger.Interface
}

func (u *Archive) Do(in Input) error {
	if err := u.FileSystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return fmt.Errorf("error while creating the output directory: %w", err)
	}

	u.Logger.Logf("Creating %s", in.OutputFilename)
	output, err := u.FileSystem.Create(in.OutputFilename)
	if err != nil {
		return fmt.Errorf("error while creating the file %s: %w", in.OutputFilename, err)
	}
	defer output.Close()

	w := zip.NewWriter(output)
	for _, e := range in.Entries {
		if err := u.addEntry(w, e); err != nil {
			return fmt.Errorf("error while adding the file %s: %w", e.InputFilename, err)
		}
	}
	if err := w.Close(); err != nil {
		return fmt.Errorf("error while writing to the file %s: %w", in.OutputFilename, err)
	}
	return nil
}

func (u *Archive) addEntry(zipWriter *zip.Writer, e Entry) error {
	stat, err := u.FileSystem.Stat(e.InputFilename)
	if err != nil {
		return fmt.Errorf("error while getting mode of the file %s: %w", e.InputFilename, err)
	}
	h := &zip.FileHeader{
		Name:     e.Filename,
		Method:   zip.Deflate,
		Modified: stat.ModTime(),
	}
	h.SetMode(stat.Mode())
	w, err := zipWriter.CreateHeader(h)
	if err != nil {
		return fmt.Errorf("error while creating a header for the file %s: %w", e.Filename, err)
	}
	input, err := u.FileSystem.Open(e.InputFilename)
	if err != nil {
		return fmt.Errorf("error while opening the file %s: %w", e.InputFilename, err)
	}
	defer input.Close()
	if _, err := io.Copy(w, input); err != nil {
		return fmt.Errorf("error while compressing the file %s: %w", e.InputFilename, err)
	}
	return nil
}
