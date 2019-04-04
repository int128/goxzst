package adaptors

import (
	"flag"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

const usage = `%[1]s %[2]s
A command for cross-build, zip archive, sha digest for each GOOS/GOARCH and template rendering.

Examples:
  To make cross-build, zip and sha256 for the default platforms:
    %[1]s

  You can set the target platforms:
    %[1]s -osarch "linux_amd64 linux_arm"

  You can pass extra arguments to go build:
    %[1]s -- -ldflags "-X main.version=$VERSION"

  You can add extra files to zip:
    %[1]s -i "LICENSE README.md"

Usage:
  %[1]s [-d DIR] [-o NAME] [-osarch "GOOS_GOARCH ..."] [-i "FILE ..."] [-t "FILE ..."] [--] [build args]

Options:
`

func NewCmd(i Cmd) adaptors.Cmd {
	return &i
}

type Cmd struct {
	dig.In
	Make   usecases.Make
	Env    adaptors.Env
	Logger adaptors.Logger
}

func (cmd *Cmd) Run(args []string, version string) int {
	wd, _ := cmd.Env.Getwd()
	var o cmdOptions
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Usage = func() {
		_, _ = fmt.Fprintf(f.Output(), usage, f.Name(), version)
		f.PrintDefaults()
	}
	f.StringVar(&o.outputDir, "d", "dist", "Output directory")
	f.StringVar(&o.outputName, "o", filepath.Base(wd), "Output name")
	f.StringVar(&o.osarch, "osarch", "linux_amd64 darwin_amd64 windows_amd64", "List of GOOS_GOARCH separated by space")
	f.StringVar(&o.archiveExtraFilenames, "i", "", "List of extra files to add to the zip, separated by space")
	f.StringVar(&o.templateFilenames, "t", "", "List of template files separated by space")
	if err := f.Parse(args[1:]); err != nil {
		return 1
	}
	platforms, err := o.platformList()
	if err != nil {
		cmd.Logger.Logf("Invalid arguments: %s", err)
		return 1
	}

	in := usecases.MakeIn{
		OutputDir:             o.outputDir,
		OutputName:            o.outputName,
		Platforms:             platforms,
		GoBuildArgs:           f.Args(),
		ArchiveExtraFilenames: o.archiveExtraFilenameList(),
		TemplateFilenames:     o.templateFilenameList(),
	}
	if err := cmd.Make.Do(in); err != nil {
		cmd.Logger.Logf("Error: %s", err)
		return 1
	}
	return 0
}

type cmdOptions struct {
	outputDir             string
	outputName            string
	osarch                string
	archiveExtraFilenames string
	templateFilenames     string
}

func (o *cmdOptions) platformList() ([]build.Platform, error) {
	var platforms []build.Platform
	for _, s := range strings.Split(o.osarch, " ") {
		p := strings.SplitN(s, "_", 2)
		if len(p) != 2 {
			return nil, errors.Errorf("osarch must be GOOS_GOARCH but was %s", s)
		}
		platforms = append(platforms, build.Platform{
			GOOS:   build.GOOS(p[0]),
			GOARCH: build.GOARCH(p[1]),
		})
	}
	return platforms, nil
}

func (o *cmdOptions) archiveExtraFilenameList() []string {
	if o.archiveExtraFilenames == "" {
		return nil
	}
	return strings.Split(o.archiveExtraFilenames, " ")
}

func (o *cmdOptions) templateFilenameList() []string {
	if o.templateFilenames == "" {
		return nil
	}
	return strings.Split(o.templateFilenames, " ")
}
