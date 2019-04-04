package usecases

import (
	"fmt"
	"path/filepath"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/models/build"
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
	Archive        usecases.Archive
	Digest         usecases.Digest
	RenderTemplate usecases.RenderTemplate
	Filesystem     adaptors.Filesystem
	Logger         adaptors.Logger
}

func (u *Make) Do(in usecases.MakeIn) error {
	var executableFilenames []string
	templateVariables := make(map[string]string)

	for _, platform := range in.Platforms {
		out, err := u.build(in, platform)
		if err != nil {
			return errors.Wrapf(err, "error while build for the platform %s", platform)
		}
		executableFilenames = append(executableFilenames, out.executableFilename)
		templateVariables[fmt.Sprintf("%s_%s_zip_sha256", platform.GOOS, platform.GOARCH)] = out.digestOut.SHA256
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

	for _, executableFilename := range executableFilenames {
		u.Logger.Logf("Removing %s", executableFilename)
		if err := u.Filesystem.Remove(executableFilename); err != nil {
			return errors.Wrapf(err, "error while removing %s", executableFilename)
		}
	}
	return nil
}

type buildOut struct {
	executableFilename string
	digestOut          *usecases.DigestOut
}

func (u *Make) build(in usecases.MakeIn, platform build.Platform) (*buildOut, error) {
	basename := filepath.Join(in.OutputDir, fmt.Sprintf("%s_%s_%s", in.OutputName, platform.GOOS, platform.GOARCH))
	executableFilename := basename
	zipFilename := executableFilename + ".zip"
	shaFilename := executableFilename + ".zip.sha256"

	if err := u.CrossBuild.Do(usecases.CrossBuildIn{
		OutputFilename: executableFilename,
		GoBuildArgs:    in.GoBuildArgs,
		Platform:       platform,
	}); err != nil {
		return nil, errors.Wrapf(err, "error while cross build")
	}

	archiveEntries := []usecases.ArchiveEntry{{
		Path:          in.OutputName,
		InputFilename: executableFilename,
	}}
	for _, f := range in.ArchiveExtraFilenames {
		archiveEntries = append(archiveEntries, usecases.ArchiveEntry{Path: f, InputFilename: f})
	}
	if err := u.Archive.Do(usecases.ArchiveIn{
		OutputFilename: zipFilename,
		Entries:        archiveEntries,
	}); err != nil {
		return nil, errors.Wrapf(err, "error while creating zip")
	}

	digestOut, err := u.Digest.Do(usecases.DigestIn{
		InputFilename:  zipFilename,
		OutputFilename: shaFilename,
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error while creating digest")
	}

	return &buildOut{
		executableFilename: executableFilename,
		digestOut:          digestOut,
	}, nil
}
