package makeall

import (
	"fmt"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/archive"
	build2 "github.com/int128/goxzst/usecases/crossbuild"
	digest2 "github.com/int128/goxzst/usecases/digest"
	"github.com/int128/goxzst/usecases/rendertemplate"
	"github.com/pkg/errors"
)

var Set = wire.NewSet(
	wire.Struct(new(MakeAll), "*"),
	wire.Bind(new(Interface), new(*MakeAll)),
)

//go:generate mockgen -destination mock_makeall/mock_makeall.go github.com/int128/goxzst/usecases/makeall Interface

type Interface interface {
	Do(in Input) error
}

type Input struct {
	OutputDir             string // optional
	OutputName            string
	Platforms             []build.Platform
	GoBuildArgs           []string
	ArchiveExtraFilenames []string
	DigestAlgorithm       *digest.Algorithm
	TemplateFilenames     []string
}

type MakeAll struct {
	CrossBuild     build2.Interface
	Archive        archive.Interface
	Digest         digest2.Interface
	RenderTemplate rendertemplate.Interface
	FileSystem     fs.Interface
	Logger         logger.Interface
}

func (u *MakeAll) Do(in Input) error {
	if in.DigestAlgorithm == nil {
		return errors.New("DigestAlgorithm must be non-nil")
	}

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
		templateVariables[fmt.Sprintf("%s_%s_executable", buildOut.platform.GOOS, buildOut.platform.GOARCH)] =
			buildOut.executableFile.Name()
		templateVariables[fmt.Sprintf("%s_%s_archive", buildOut.platform.GOOS, buildOut.platform.GOARCH)] =
			buildOut.archiveFile.Name()
		templateVariables[fmt.Sprintf("%s_%s_digest", buildOut.platform.GOOS, buildOut.platform.GOARCH)] =
			buildOut.digestFile.Name()
	}
	for _, t := range in.TemplateFilenames {
		if err := u.RenderTemplate.Do(rendertemplate.Input{
			InputFilename:  t,
			OutputFilename: filepath.Join(in.OutputDir, filepath.Base(t)),
			Variables:      templateVariables,
		}); err != nil {
			return errors.Wrapf(err, "error while rendering templates")
		}
	}

	for _, buildOut := range buildOuts {
		name := buildOut.executableFile.Name()
		u.Logger.Logf("Removing %s", name)
		if err := u.FileSystem.Remove(name); err != nil {
			return errors.Wrapf(err, "error while removing %s", name)
		}
	}
	return nil
}

type buildOut struct {
	platform       build.Platform
	executableFile build.ExecutableFile
	archiveFile    build.ArchiveFile
	digestFile     build.DigestFile
}

func (u *MakeAll) build(in Input, platform build.Platform) (*buildOut, error) {
	basename := filepath.Join(in.OutputDir, fmt.Sprintf("%s_%s_%s", in.OutputName, platform.GOOS, platform.GOARCH))

	builtExecutableFile := build.ExecutableFile{
		Base:     basename,
		Platform: platform,
	}
	if err := u.CrossBuild.Do(build2.Input{
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
	archiveEntries := []archive.Entry{{
		Filename:      executableInArchive.Name(),
		InputFilename: builtExecutableFile.Name(),
	}}
	for _, f := range in.ArchiveExtraFilenames {
		archiveEntries = append(archiveEntries, archive.Entry{
			Filename:      f,
			InputFilename: f,
		})
	}
	if err := u.Archive.Do(archive.Input{
		OutputFilename: archiveFile.Name(),
		Entries:        archiveEntries,
	}); err != nil {
		return nil, errors.Wrapf(err, "error while archive")
	}

	digestFile := build.DigestFile{
		Base:   archiveFile.Name(),
		Suffix: in.DigestAlgorithm.Suffix,
	}
	if err := u.Digest.Do(digest2.Input{
		InputFilename:  archiveFile.Name(),
		OutputFilename: digestFile.Name(),
		Algorithm:      in.DigestAlgorithm,
	}); err != nil {
		return nil, errors.Wrapf(err, "error while digest")
	}

	return &buildOut{
		platform:       platform,
		executableFile: builtExecutableFile,
		archiveFile:    archiveFile,
		digestFile:     digestFile,
	}, nil
}