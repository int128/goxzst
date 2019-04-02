package usecases

import (
	"archive/zip"
	"io"
	"os"

	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
)

func NewCreateZip() usecases.CreateZip {
	return &CreateZip{}
}

type CreateZip struct{}

func (u *CreateZip) Do(in usecases.CreateZipIn) error {
	output, err := os.Create(in.OutputFilename)
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

func (*CreateZip) addEntry(zipWriter *zip.Writer, e usecases.ZipEntry) error {
	h := &zip.FileHeader{
		Name:   e.Path,
		Method: zip.Deflate,
	}
	stat, err := os.Stat(e.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while getting status of the file %s", e.InputFilename)
	}
	h.SetMode(stat.Mode())
	w, err := zipWriter.CreateHeader(h)
	if err != nil {
		return errors.Wrapf(err, "error while creating a header for the file %s", e.Path)
	}
	input, err := os.Open(e.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while opening the file %s", e.InputFilename)
	}
	defer input.Close()
	if _, err := io.Copy(w, input); err != nil {
		return errors.Wrapf(err, "error while compressing the file %s", e.InputFilename)
	}
	return nil
}
