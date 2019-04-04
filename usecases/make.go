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
	FileSystem     adaptors.FileSystem
	Logger         adaptors.Logger
}

func (u *Make) Do(in usecases.MakeIn) error {
	var buildOuts []*buildOut
	for _, platform := range in.Platforms {
		out, err := u.build(in, platform)
		if err != nil {
			return errors.Wrapf(err, "error while build for the platform %s", platform)
		}
		buildOuts = append(buildOuts, out)
	}

	templateVariables := make(map[string]string)
	for _, buildOut := range buildOuts {
		templateVariables[fmt.Sprintf("%s_%s_zip_sha256", buildOut.platform.GOOS, buildOut.platform.GOARCH)] =
			buildOut.digestOut.SHA256
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

	for _, buildOut := range buildOuts {
		u.Logger.Logf("Removing %s", buildOut.executableFilename)
		if err := u.FileSystem.Remove(buildOut.executableFilename); err != nil {
			return errors.Wrapf(err, "error while removing %s", buildOut.executableFilename)
		}
	}
	return nil
}

type buildOut struct {
	platform           build.Platform
	executableFilename string
	digestOut          *usecases.DigestOut
}

func (u *Make) build(in usecases.MakeIn, platform build.Platform) (*buildOut, error) {
	basename := filepath.Join(in.OutputDir, fmt.Sprintf("%s_%s_%s", in.OutputName, platform.GOOS, platform.GOARCH))

	builtExecutableFile := build.ExecutableFile{
		Base:     basename,
		Platform: platform,
	}
	if err := u.CrossBuild.Do(usecases.CrossBuildIn{
		OutputFilename: builtExecutableFile.Name(),
		GoBuildArgs:    in.GoBuildArgs,
		Platform:       platform,
	}); err != nil {
		return nil, errors.Wrapf(err, "error while cross build")
	}

	archiveFile := build.ArchiveFile{
		Base:   basename,
		Suffix: ".zip",
	}
	executableInArchive := build.ExecutableFile{
		Base:     in.OutputName,
		Platform: platform,
	}
	archiveEntries := []usecases.ArchiveEntry{{
		Filename:      executableInArchive.Name(),
		InputFilename: builtExecutableFile.Name(),
	}}
	for _, f := range in.ArchiveExtraFilenames {
		archiveEntries = append(archiveEntries, usecases.ArchiveEntry{
			Filename:      f,
			InputFilename: f,
		})
	}
	if err := u.Archive.Do(usecases.ArchiveIn{
		OutputFilename: archiveFile.Name(),
		Entries:        archiveEntries,
	}); err != nil {
		return nil, errors.Wrapf(err, "error while archive")
	}

	digestFile := build.DigestFile{
		Base:   archiveFile.Name(),
		Suffix: ".sha256",
	}
	digestOut, err := u.Digest.Do(usecases.DigestIn{
		InputFilename:  archiveFile.Name(),
		OutputFilename: digestFile.Name(),
	})
	if err != nil {
		return nil, errors.Wrapf(err, "error while digest")
	}

	return &buildOut{
		platform:           platform,
		executableFilename: builtExecutableFile.Name(),
		digestOut:          digestOut,
	}, nil
}
