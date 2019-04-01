package main

import (
	"os"
	"text/template"

	"github.com/pkg/errors"
)

type RenderTemplateIn struct {
	Entries []TemplateEntry
	Data    map[string]string
}

type TemplateEntry struct {
	InputFilename  string
	OutputFilename string
}

type RenderTemplate struct{}

func (u *RenderTemplate) Do(in RenderTemplateIn) error {
	for _, e := range in.Entries {
		if err := u.render(e, in.Data); err != nil {
			return errors.Wrapf(err, "error while rendering the template")
		}
	}
	return nil
}

func (*RenderTemplate) render(e TemplateEntry, data map[string]string) error {
	t, err := template.ParseFiles(e.InputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while loading templates")
	}
	t.Option("missingkey=zero")

	output, err := os.Create(e.OutputFilename)
	if err != nil {
		return errors.Wrapf(err, "error while creating the file %s", e.OutputFilename)
	}
	defer output.Close()

	if err := t.Execute(output, data); err != nil {
		return errors.Wrapf(err, "error while rendering the template %s", e.InputFilename)
	}
	return nil
}
