package usecases

import "github.com/int128/goxzst/models/build"

type Make interface {
	Do(in MakeIn) error
}

type MakeIn struct {
	OutputDir         string // optional
	OutputName        string
	Targets           []build.Target
	GoBuildArgs       []string
	TemplateFilenames []string
	TemplateVariables map[string]string
}

type CrossBuild interface {
	Do(in CrossBuildIn) error
}

type CrossBuildIn struct {
	OutputFilename string
	GoBuildArgs    []string
	Target         build.Target
}

type CreateZip interface {
	Do(in CreateZipIn) error
}

type CreateZipIn struct {
	OutputFilename string
	Entries        []ZipEntry
}

type ZipEntry struct {
	Path          string
	InputFilename string
}

type CreateSHA interface {
	Do(in CreateSHAIn) (*CreateSHAOut, error)
}

type CreateSHAIn struct {
	InputFilename  string
	OutputFilename string
}

type CreateSHAOut struct {
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
