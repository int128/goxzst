package usecases

import (
	"fmt"
	"path/filepath"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

func NewMake(i Make) usecases.Make {
	return &i
}

type Make struct {
	dig.In
	CrossBuild     usecases.CrossBuild
	CreateZip      usecases.CreateZip
	CreateSHA      usecases.CreateSHA
	RenderTemplate usecases.RenderTemplate
	Filesystem     adaptors.Filesystem
}

func (u *Make) Do(in usecases.MakeIn) error {
	if in.OutputDir != "" {
		if err := u.Filesystem.MkdirAll(in.OutputDir); err != nil {
			return errors.Wrapf(err, "error while creating the output directory")
		}
	}

	templateVariables := make(map[string]string)
	for _, platform := range in.Platforms {
		executableFilename := filepath.Join(in.OutputDir,
			fmt.Sprintf("%s_%s_%s", in.OutputName, platform.GOOS, platform.GOARCH))
		zipFilename := executableFilename + ".zip"
		shaFilename := executableFilename + ".zip.sha256"

		if err := u.CrossBuild.Do(usecases.CrossBuildIn{
			OutputFilename: executableFilename,
			GoBuildArgs:    in.GoBuildArgs,
			Platform:       platform,
		}); err != nil {
			return errors.Wrapf(err, "error while cross build")
		}

		if err := u.CreateZip.Do(usecases.CreateZipIn{
			OutputFilename: zipFilename,
			Entries: []usecases.ZipEntry{
				{
					Path:          in.OutputName,
					InputFilename: executableFilename,
				},
			},
		}); err != nil {
			return errors.Wrapf(err, "error while creating zip")
		}

		shaOut, err := u.CreateSHA.Do(usecases.CreateSHAIn{
			InputFilename:  zipFilename,
			OutputFilename: shaFilename,
		})
		if err != nil {
			return errors.Wrapf(err, "error while creating digest")
		}

		templateVariables[fmt.Sprintf("%s_%s_zip_sha256", platform.GOOS, platform.GOARCH)] = shaOut.SHA256
	}

	for _, t := range in.TemplateFilenames {
		if err := u.RenderTemplate.Do(usecases.RenderTemplateIn{
			InputFilename:  t,
			OutputFilename: filepath.Join(in.OutputDir, filepath.Base(t)),
			Variables:      templateVariables,
		}); err != nil {
			return errors.Wrapf(err, "error while rendering templates")
		}
	}
	return nil
}
