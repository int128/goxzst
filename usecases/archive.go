package usecases

import (
	"archive/zip"
	"io"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewArchive(i Archive) usecases.Archive {
	return &i
}

type Archive struct {
	dig.In
	Filesystem adaptors.Filesystem
}

func (u *Archive) Do(in usecases.ArchiveIn) error {
	output, err := u.Filesystem.Create(in.OutputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while creating the file %s", in.OutputFilename)
	}
	defer output.Close()

	w := zip.NewWriter(output)
	for _, e := range in.Entries {
		if err := u.addEntry(w, e); err != nil {
			return errors.Wrapf(err, "error while adding the file %s", e.InputFilename)
		}
	}
	if err := w.Close(); err != nil {
		return errors.Wrapf(err, "error while writing to the file %s", in.OutputFilename)
	}
	return nil
}

func (u *Archive) addEntry(zipWriter *zip.Writer, e usecases.ArchiveEntry) error {
	h := &zip.FileHeader{
		Name:   e.Path,
		Method: zip.Deflate,
	}
	mode, err := u.Filesystem.GetMode(e.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while getting mode of the file %s", e.InputFilename)
	}
	h.SetMode(mode)
	w, err := zipWriter.CreateHeader(h)
	if err != nil {
		return errors.Wrapf(err, "error while creating a header for the file %s", e.Path)
	}
	input, err := u.Filesystem.Open(e.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while opening the file %s", e.InputFilename)
	}
	defer input.Close()
	if _, err := io.Copy(w, input); err != nil {
		return errors.Wrapf(err, "error while compressing the file %s", e.InputFilename)
	}
	return nil
}
