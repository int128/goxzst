package usecases

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
)

type RenderTemplate struct {
	Env        adaptors.Env
	FileSystem adaptors.FileSystem
	Logger     adaptors.Logger
}

func (u *RenderTemplate) Do(in usecases.RenderTemplateIn) error {
	if err := u.FileSystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return errors.Wrapf(err, "error while creating the output directory")
	}

	u.Logger.Logf("Creating %s from the template %s", in.OutputFilename, in.InputFilename)
	tpl, err := template.New(filepath.Base(in.InputFilename)).
		Funcs(template.FuncMap{
			"env": u.env,
			"sha256": func(filename string) (string, error) {
				return u.digest(filename, digest.SHA256)
			},
			"sha512": func(filename string) (string, error) {
				return u.digest(filename, digest.SHA512)
			},
		}).
		ParseFiles(in.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while loading templates")
	}

	output, err := u.FileSystem.Create(in.OutputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while creating the file %s", in.OutputFilename)
	}
	defer output.Close()

	if err := tpl.Execute(output, in.Variables); err != nil {
		return errors.Wrapf(err, "error while rendering the template %s", in.InputFilename)
	}
	return nil
}

func (u *RenderTemplate) env(key string) (string, error) {
	value, ok := u.Env.LookupEnv(key)
	if !ok {
		return "", errors.Errorf("no such environment variable %s", key)
	}
	return value, nil
}

func (u *RenderTemplate) digest(filename string, algorithm *digest.Algorithm) (string, error) {
	r, err := u.FileSystem.Open(filename)
	if err != nil {
		return "", errors.Errorf("error while opening %s", filename)
	}
	defer r.Close()
	w := algorithm.NewHash()
	if _, err := io.Copy(w, r); err != nil {
		return "", errors.Wrapf(err, "error while computing digest of the file %s", filename)
	}
	h := fmt.Sprintf("%x", w.Sum(nil))
	return h, nil
}
