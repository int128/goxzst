package usecases

import (
	"path/filepath"
	"text/template"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewRenderTemplate(i RenderTemplate) usecases.RenderTemplate {
	return &i
}

type RenderTemplate struct {
	dig.In
	Env        adaptors.Env
	Filesystem adaptors.Filesystem
}

func (u *RenderTemplate) Do(in usecases.RenderTemplateIn) error {
	tpl, err := template.New(filepath.Base(in.InputFilename)).
		Funcs(template.FuncMap{
			"env": u.env,
		}).
		ParseFiles(in.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while loading templates")
	}

	output, err := u.Filesystem.Create(in.OutputFilename)
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
