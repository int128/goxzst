package adaptors

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
	"go.uber.org/dig"
)

const usage = `Crossbuild, zip, shasum for each GOOS/GOARCH and render templates.

Usage:
  %[1]s [-d DIR] [-o NAME] [-osarch "GOOS_GOARCH ..."] [-t "FILE ..."] [-tvar "KEY=VALUE ..."] [--] [build args]

Options:
`

func NewCmd(i Cmd) adaptors.Cmd {
	return &i
}

type Cmd struct {
	dig.In
	Make   usecases.Make
	Logger adaptors.Logger
}

func (cmd *Cmd) Run(args []string) int {
	wd, _ := os.Getwd()
	var o cmdOptions
	f := flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	f.Usage = func() {
		_, _ = fmt.Fprintf(f.Output(), usage, f.Name())
		f.PrintDefaults()
	}
	f.StringVar(&o.outputDir, "d", "dist", "Output directory")
	f.StringVar(&o.outputName, "o", filepath.Base(wd), "Output name")
	f.StringVar(&o.osarch, "osarch", "linux_amd64 darwin_amd64 windows_amd64", "List of GOOS_GOARCH separated by space")
	f.StringVar(&o.templateFilenames, "t", "", "List of template files separated by space")
	f.StringVar(&o.templateVariables, "tvar", "", "List of template variables as KEY=VALUE separated by space")
	if err := f.Parse(os.Args[1:]); err != nil {
		return 1
	}
	targets, err := o.targetList()
	if err != nil {
		cmd.Logger.Logf("Invalid arguments: %s", err)
		return 1
	}
	templateVariableList, err := o.templateVariableMap()
	if err != nil {
		cmd.Logger.Logf("Invalid arguments: %s", err)
		return 1
	}

	in := usecases.MakeIn{
		OutputDir:         o.outputDir,
		OutputName:        o.outputName,
		Targets:           targets,
		GoBuildArgs:       f.Args(),
		TemplateFilenames: o.templateFilenameList(),
		TemplateVariables: templateVariableList,
	}
	if err := cmd.Make.Do(in); err != nil {
		cmd.Logger.Logf("Error: %s", err)
		return 1
	}
	return 0
}

type cmdOptions struct {
	outputDir         string
	outputName        string
	osarch            string
	templateFilenames string
	templateVariables string
}

func (o *cmdOptions) targetList() ([]build.Target, error) {
	var targets []build.Target
	for _, s := range strings.Split(o.osarch, " ") {
		p := strings.SplitN(s, "_", 2)
		if len(p) != 2 {
			return nil, errors.Errorf("osarch must be GOOS_GOARCH but was %s", s)
		}
		targets = append(targets, build.Target{
			GOOS:   build.GOOS(p[0]),
			GOARCH: build.GOARCH(p[1]),
		})
	}
	return targets, nil
}

func (o *cmdOptions) templateFilenameList() []string {
	if o.templateFilenames == "" {
		return []string{}
	}
	return strings.Split(o.templateFilenames, " ")
}

func (o *cmdOptions) templateVariableMap() (map[string]string, error) {
	vars := make(map[string]string)
	if o.templateVariables == "" {
		return vars, nil
	}
	for _, s := range strings.Split(o.templateVariables, " ") {
		p := strings.SplitN(s, "=", 2)
		if len(p) != 2 {
			return nil, errors.Errorf("template variable must be KEY=VALUE but was %s", s)
		}
		k, v := p[0], p[1]
		vars[k] = v
	}
	return vars, nil
}
