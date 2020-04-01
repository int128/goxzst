// Package xzs provides the use-case to make an executable, archive and digest for a platform.
package xzs

import (
	"fmt"
	"path/filepath"

	"github.com/google/wire"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/archive"
	"github.com/int128/goxzst/usecases/crossbuild"
	digestUseCase "github.com/int128/goxzst/usecases/digest"
)

var Set = wire.NewSet(
	wire.Struct(new(XZS), "*"),
	wire.Bind(new(Interface), new(*XZS)),
)

//go:generate mockgen -destination mock_xzs/mock_xzs.go github.com/int128/goxzst/usecases/xzs Interface

type Interface interface {
	Do(in Input) (*build.Artifact, error)
}

type Input struct {
	OutputDir             string // optional
	OutputName            string
	Platform              build.Platform
	GoBuildArgs           []string
	ArchiveExtraFilenames []string
	DigestAlgorithm       *digest.Algorithm
}

type XZS struct {
	CrossBuild crossbuild.Interface
	Archive    archive.Interface
	Digest     digestUseCase.Interface
}

func (u *XZS) Do(in Input) (*build.Artifact, error) {
	basename := filepath.Join(in.OutputDir, fmt.Sprintf("%s_%s_%s", in.OutputName, in.Platform.GOOS, in.Platform.GOARCH))

	builtExecutableFile := build.ExecutableFile{
		Base:     basename,
		Platform: in.Platform,
	}
	if err := u.CrossBuild.Do(crossbuild.Input{
		OutputFilename: builtExecutableFile.Name(),
		GoBuildArgs:    in.GoBuildArgs,
		Platform:       in.Platform,
	}); err != nil {
		return nil, fmt.Errorf("error while cross build: %w", err)
	}

	archiveFile := build.ArchiveFile{
		Base:   basename,
		Suffix: ".zip",
	}
	executableInArchive := build.ExecutableFile{
		Base:     in.OutputName,
		Platform: in.Platform,
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
		return nil, fmt.Errorf("error while archive: %w", err)
	}

	digestFile := build.DigestFile{
		Base:   archiveFile.Name(),
		Suffix: in.DigestAlgorithm.Suffix,
	}
	if err := u.Digest.Do(digestUseCase.Input{
		InputFilename:  archiveFile.Name(),
		OutputFilename: digestFile.Name(),
		Algorithm:      in.DigestAlgorithm,
	}); err != nil {
		return nil, fmt.Errorf("error while digest: %w", err)
	}

	return &build.Artifact{
		Platform:       in.Platform,
		ExecutableFile: builtExecutableFile,
		ArchiveFile:    archiveFile,
		DigestFile:     digestFile,
	}, nil
}
