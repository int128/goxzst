package usecases

import "github.com/int128/goxzst/models/build"

//go:generate mockgen -package mock_usecases -destination ../mock_usecases/mock_usecases.go github.com/int128/goxzst/usecases/interfaces Make,CrossBuild,Archive,Digest,RenderTemplate

type Make interface {
	Do(in MakeIn) error
}

type MakeIn struct {
	OutputDir         string // optional
	OutputName        string
	Platforms         []build.Platform
	GoBuildArgs       []string
	TemplateFilenames []string
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
	Path          string
	InputFilename string
}

type Digest interface {
	Do(in DigestIn) (*DigestOut, error)
}

type DigestIn struct {
	InputFilename  string
	OutputFilename string
}

type DigestOut struct {
	SHA256 string
}

type RenderTemplate interface {
	Do(in RenderTemplateIn) error
}

type RenderTemplateIn struct {
	InputFilename  string
	OutputFilename string
	Variables      map[string]string
}
