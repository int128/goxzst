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
		Option("missingkey=zero").
		Funcs(template.FuncMap{
			"env": func(key string) string {
				return u.Env.Getenv(key)
			},
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
