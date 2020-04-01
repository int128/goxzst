// Package cmd provides the command line interface (CLI).
package cmd

import (
	"flag"
	"fmt"
	"strings"

	"github.com/google/wire"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/logger"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/xzst"
)

const usage = `%[1]s %[2]s
A command for cross-build, zip archive, sha digest for each GOOS/GOARCH and template rendering.

Examples:
  To make cross-build, zip and sha256 for the default platforms:
    %[1]s -o NAME

  You can set the target platforms:
    %[1]s -o NAME -osarch "linux_amd64 linux_arm"

  You can pass extra arguments to go build:
    %[1]s -o NAME -- -ldflags "-X main.version=$VERSION"

  You can add extra files to zip:
    %[1]s -o NAME -i "LICENSE README.md"

Usage:
  %[1]s -o NAME [-d DIR] [-osarch "GOOS_GOARCH ..."] [-i "FILE ..."] [-a ALGORITHM] [-t "FILE ..."] [--] [build args]

Options:
`

var Set = wire.NewSet(
	wire.Struct(new(Cmd), "*"),
	wire.Bind(new(Interface), new(*Cmd)),
)

type Interface interface {
	Run(args []string, version string) int
}

type Cmd struct {
	XZST   xzst.Interface
	Env    env.Interface
	Logger logger.Interface
}

// Run parses the command line arguments and executes the corresponding use-case.
func (cmd *Cmd) Run(args []string, version string) int {
	var o cmdOptions
	f := flag.NewFlagSet(args[0], flag.ExitOnError)
	f.Usage = func() {
		_, _ = fmt.Fprintf(f.Output(), usage, f.Name(), version)
		f.PrintDefaults()
	}
	f.StringVar(&o.outputName, "o", "", "Output name (mandatory)")
	f.StringVar(&o.outputDir, "d", "dist", "Output directory")
	f.StringVar(&o.osarch, "osarch", "linux_amd64 darwin_amd64 windows_amd64", "List of GOOS_GOARCH separated by space")
	f.StringVar(&o.archiveExtraFilenames, "i", "", "List of extra files to add to the zip, separated by space")
	f.StringVar(&o.digestAlgorithm, "a", "sha256", fmt.Sprintf("Digest algorithm. One of (%s)", availableDigestAlgorithms()))
	f.StringVar(&o.templateFilenames, "t", "", "List of template files separated by space")
	if err := f.Parse(args[1:]); err != nil {
		return 1
	}
	if o.outputName == "" {
		cmd.Logger.Logf("You need to set output name by -o option")
		return 1
	}
	platforms, err := o.platformList()
	if err != nil {
		cmd.Logger.Logf("Invalid arguments: %s", err)
		return 1
	}
	digestAlgorithm, err := digest.NewAlgorithm(o.digestAlgorithm)
	if err != nil {
		cmd.Logger.Logf("Invalid digest algorithm: %s", err)
		return 1
	}

	in := xzst.Input{
		OutputDir:             o.outputDir,
		OutputName:            o.outputName,
		Platforms:             platforms,
		GoBuildArgs:           f.Args(),
		ArchiveExtraFilenames: o.archiveExtraFilenameList(),
		DigestAlgorithm:       digestAlgorithm,
		TemplateFilenames:     o.templateFilenameList(),
	}
	if err := cmd.XZST.Do(in); err != nil {
		cmd.Logger.Logf("Error: %s", err)
		return 1
	}
	return 0
}

func availableDigestAlgorithms() string {
	var names []string
	for _, alg := range digest.AvailableAlgorithms {
		names = append(names, alg.Name)
	}
	return strings.Join(names, "|")
}

type cmdOptions struct {
	outputName            string
	outputDir             string
	osarch                string
	archiveExtraFilenames string
	digestAlgorithm       string
	templateFilenames     string
}

func (o *cmdOptions) platformList() ([]build.Platform, error) {
	var platforms []build.Platform
	for _, s := range strings.Split(o.osarch, " ") {
		p := strings.SplitN(s, "_", 2)
		if len(p) != 2 {
			return nil, fmt.Errorf("osarch must be GOOS_GOARCH but was %s", s)
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
