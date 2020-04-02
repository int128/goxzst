package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/env/mock_env"
	"github.com/int128/goxzst/adaptors/log"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/xzst"
	"github.com/int128/goxzst/usecases/xzst/mock_xzst"
)

func TestCmd_Run(t *testing.T) {
	log.Printf = t.Logf
	const version = "dummyVersionString"
	defaultPlatforms := []build.Platform{
		{GOOS: "linux", GOARCH: "amd64"},
		{GOOS: "darwin", GOARCH: "amd64"},
		{GOOS: "windows", GOARCH: "amd64"},
	}

	t.Run("NoArg", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		cmd := Cmd{
			XZST: mock_xzst.NewMockInterface(ctrl),
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst"}, version)
		if exitCode != 1 {
			t.Errorf("exitCode wants 1 but %d", exitCode)
		}
	})

	t.Run("MinimumArgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA256,
				Parallelism:     defaultParallelism,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithGoBuildArgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{"-ldflags", "-X foo=bar"},
				DigestAlgorithm: digest.SHA256,
				Parallelism:     defaultParallelism,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "--", "-ldflags", "-X foo=bar"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithPlatforms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:  "dist",
				OutputName: "hello",
				Platforms: []build.Platform{
					{GOOS: "linux", GOARCH: "arm"},
				},
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA256,
				Parallelism:     defaultParallelism,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-osarch", "linux_arm"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithExtraFilesToZip", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:             "dist",
				OutputName:            "hello",
				Platforms:             defaultPlatforms,
				GoBuildArgs:           []string{},
				ArchiveExtraFilenames: []string{"README.md", "LICENSE"},
				DigestAlgorithm:       digest.SHA256,
				Parallelism:           defaultParallelism,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-i", "README.md LICENSE"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithDigestAlgorithm", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA512,
				Parallelism:     defaultParallelism,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-a", "sha512"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithTemplates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:         "dist",
				OutputName:        "hello",
				Platforms:         defaultPlatforms,
				GoBuildArgs:       []string{},
				DigestAlgorithm:   digest.SHA256,
				TemplateFilenames: []string{"template1", "template2"},
				Parallelism:       defaultParallelism,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-t", "template1 template2"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithParallelism", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		xzstUseCase := mock_xzst.NewMockInterface(ctrl)
		xzstUseCase.EXPECT().
			Do(xzst.Input{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA256,
				Parallelism:     98,
			})

		cmd := Cmd{
			XZST: xzstUseCase,
			Env:  mock_env.NewMockInterface(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-parallelism", "98"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})
}
