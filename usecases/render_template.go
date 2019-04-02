package usecases

import (
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
	Filesystem adaptors.Filesystem
}

func (u *RenderTemplate) Do(in usecases.RenderTemplateIn) error {
	tpl, err := template.ParseFiles(in.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while loading templates")
	}
	tpl.Option("missingkey=zero")

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
