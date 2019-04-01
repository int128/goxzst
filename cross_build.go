package main

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

type CrossBuildIn struct {
	OutputFilename string
	Args           []string
	GOOS           string
	GOARCH         string
}

type CrossBuild struct{}

func (*CrossBuild) Do(in CrossBuildIn) error {
	cmd := exec.Command("go", "build", "-o", in.OutputFilename)
	cmd.Args = append(cmd.Args, in.Args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(), "GOOS="+in.GOOS, "GOARCH="+in.GOARCH)
	if err := cmd.Run(); err != nil {
		return errors.Wrapf(err, "go build error")
	}
	return nil
}
