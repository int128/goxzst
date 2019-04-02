package usecases

import (
	"crypto/sha256"
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
)

func NewCreateSHA() usecases.CreateSHA {
	return &CreateSHA{}
}

type CreateSHA struct{}

func (*CreateSHA) Do(in usecases.CreateSHAIn) (*usecases.CreateSHAOut, error) {
	input, err := os.Open(in.InputFilename)
	if err != nil {
		return nil, errors.Wrapf(err, "error while opening the file %s", in.InputFilename)
	}
	defer input.Close()

	w := sha256.New()
	if _, err := io.Copy(w, input); err != nil {
		return nil, errors.Wrapf(err, "error while computing digest of the file %s", in.InputFilename)
	}
	sum := w.Sum(nil)
	h := fmt.Sprintf("%x", sum)

	if err := ioutil.WriteFile(in.OutputFilename, []byte(h), 0644); err != nil {
		return nil, errors.Wrapf(err, "error while writing to the file %s", in.OutputFilename)
	}
	return &usecases.CreateSHAOut{
		SHA256: h,
	}, nil
}
