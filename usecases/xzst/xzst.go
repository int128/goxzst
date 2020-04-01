// Package xzst provides the use-case to make the archives, digests and templates (XZST).
package xzst

import (
	"errors"
	"fmt"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/fs"
	"github.com/int128/goxzst/adaptors/log"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/rendertemplate"
	"github.com/int128/goxzst/usecases/xzs"
)

var Set = wire.NewSet(
	wire.Struct(new(XZST), "*"),
	wire.Bind(new(Interface), new(*XZST)),
)

//go:generate mockgen -destination mock_xzst/mock_xzst.go github.com/int128/goxzst/usecases/xzst Interface

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

type XZST struct {
	XZS            xzs.Interface
	RenderTemplate rendertemplate.Interface
	FileSystem     fs.Interface
}

func (u *XZST) Do(in Input) error {
	if in.DigestAlgorithm == nil {
		return errors.New("DigestAlgorithm must be non-nil")
	}

	var artifacts []*build.Artifact
	for _, platform := range in.Platforms {
		artifact, err := u.XZS.Do(xzs.Input{
			OutputDir:             in.OutputDir,
			OutputName:            in.OutputName,
			Platform:              platform,
			GoBuildArgs:           in.GoBuildArgs,
			ArchiveExtraFilenames: in.ArchiveExtraFilenames,
			DigestAlgorithm:       in.DigestAlgorithm,
		})
		if err != nil {
			return fmt.Errorf("error while build for the platform %s: %w", platform, err)
		}
		artifacts = append(artifacts, artifact)
	}

	templateVariables := make(map[string]string)
	for _, artifact := range artifacts {
		templateVariables[fmt.Sprintf("%s_%s_executable", artifact.Platform.GOOS, artifact.Platform.GOARCH)] =
			artifact.ExecutableFile.Name()
		templateVariables[fmt.Sprintf("%s_%s_archive", artifact.Platform.GOOS, artifact.Platform.GOARCH)] =
			artifact.ArchiveFile.Name()
		templateVariables[fmt.Sprintf("%s_%s_digest", artifact.Platform.GOOS, artifact.Platform.GOARCH)] =
			artifact.DigestFile.Name()
	}
	for _, t := range in.TemplateFilenames {
		if err := u.RenderTemplate.Do(rendertemplate.Input{
			InputFilename:  t,
			OutputFilename: filepath.Join(in.OutputDir, filepath.Base(t)),
			Variables:      templateVariables,
		}); err != nil {
			return fmt.Errorf("error while rendering templates: %w", err)
		}
	}

	for _, artifact := range artifacts {
		name := artifact.ExecutableFile.Name()
		log.Printf("Removing %s", name)
		if err := u.FileSystem.Remove(name); err != nil {
			return fmt.Errorf("error while removing %s: %w", name, err)
		}
	}
	return nil
}
