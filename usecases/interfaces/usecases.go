package usecases

import (
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
)

//go:generate mockgen -package mock_usecases -destination ../mock_usecases/mock_usecases.go github.com/int128/goxzst/usecases/interfaces Make,CrossBuild,Archive,Digest,RenderTemplate

type Make interface {
	Do(in MakeIn) error
}

type MakeIn struct {
	OutputDir             string // optional
	OutputName            string
	Platforms             []build.Platform
	GoBuildArgs           []string
	ArchiveExtraFilenames []string
	DigestAlgorithm       *digest.Algorithm
	TemplateFilenames     []string
}

type CrossBuild interface {
	Do(in CrossBuildIn) error
}

type CrossBuildIn struct {
	OutputFilename string
	GoBuildArgs    []string
	Platform       build.Platform
}

type Archive interface {
	Do(in ArchiveIn) error
}

type ArchiveIn struct {
	OutputFilename string
	Entries        []ArchiveEntry
}

type ArchiveEntry struct {
	Filename      string // filename in the archive
	InputFilename string
}

type Digest interface {
	Do(in DigestIn) error
}

type DigestIn struct {
	InputFilename  string
	OutputFilename string
	Algorithm      *digest.Algorithm
}

type RenderTemplate interface {
	Do(in RenderTemplateIn) error
}

type RenderTemplateIn struct {
	InputFilename  string
	OutputFilename string
	Variables      map[string]string
}
