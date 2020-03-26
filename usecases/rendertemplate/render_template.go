package rendertemplate

import (
	"fmt"
	"io"
	"path/filepath"
	"text/template"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/models/digest"
)

var Set = wire.NewSet(
	wire.Struct(new(RenderTemplate), "*"),
	wire.Bind(new(Interface), new(*RenderTemplate)),
)

//go:generate mockgen -destination mock_rendertemplate/mock_rendertemplate.go github.com/int128/goxzst/usecases/rendertemplate Interface

type Interface interface {
	Do(in Input) error
}

type Input struct {
	InputFilename  string
	OutputFilename string
	Variables      map[string]string
}

type RenderTemplate struct {
	Env        env.Interface
	FileSystem fs.Interface
	Logger     logger.Interface
}

func (u *RenderTemplate) Do(in Input) error {
	if err := u.FileSystem.MkdirAll(filepath.Dir(in.OutputFilename)); err != nil {
		return fmt.Errorf("error while creating the output directory: %w", err)
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
		return fmt.Errorf("error while loading templates: %w", err)
	}

	output, err := u.FileSystem.Create(in.OutputFilename)
	if err != nil {
		return fmt.Errorf("error while creating the file %s: %w", in.OutputFilename, err)
	}
	defer output.Close()

	if err := tpl.Execute(output, in.Variables); err != nil {
		return fmt.Errorf("error while rendering the template %s: %w", in.InputFilename, err)
	}
	return nil
}

func (u *RenderTemplate) env(key string) (string, error) {
	value, ok := u.Env.LookupEnv(key)
	if !ok {
		return "", fmt.Errorf("no such environment variable %s", key)
	}
	return value, nil
}

func (u *RenderTemplate) digest(filename string, algorithm *digest.Algorithm) (string, error) {
	r, err := u.FileSystem.Open(filename)
	if err != nil {
		return "", fmt.Errorf("error while opening %s", filename)
	}
	defer r.Close()
	w := algorithm.NewHash()
	if _, err := io.Copy(w, r); err != nil {
		return "", fmt.Errorf("error while computing digest of the file %s: %w", filename, err)
	}
	h := fmt.Sprintf("%x", w.Sum(nil))
	return h, nil
}
