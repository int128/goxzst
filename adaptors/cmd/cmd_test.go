package cmd

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases"
	"github.com/int128/goxzst/usecases/mock_usecases"
)

func TestCmd_Run(t *testing.T) {
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
			Make:   mock_usecases.NewMockMake(ctrl),
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst"}, version)
		if exitCode != 1 {
			t.Errorf("exitCode wants 1 but %d", exitCode)
		}
	})

	t.Run("MinimumArgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA256,
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithGoBuildArgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{"-ldflags", "-X foo=bar"},
				DigestAlgorithm: digest.SHA256,
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "--", "-ldflags", "-X foo=bar"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithPlatforms", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:  "dist",
				OutputName: "hello",
				Platforms: []build.Platform{
					{GOOS: "linux", GOARCH: "arm"},
				},
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA256,
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-osarch", "linux_arm"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithExtraFilesToZip", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:             "dist",
				OutputName:            "hello",
				Platforms:             defaultPlatforms,
				GoBuildArgs:           []string{},
				ArchiveExtraFilenames: []string{"README.md", "LICENSE"},
				DigestAlgorithm:       digest.SHA256,
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-i", "README.md LICENSE"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithDigestAlgorithm", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:       "dist",
				OutputName:      "hello",
				Platforms:       defaultPlatforms,
				GoBuildArgs:     []string{},
				DigestAlgorithm: digest.SHA512,
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-a", "sha512"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})

	t.Run("WithTemplates", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:         "dist",
				OutputName:        "hello",
				Platforms:         defaultPlatforms,
				GoBuildArgs:       []string{},
				DigestAlgorithm:   digest.SHA256,
				TemplateFilenames: []string{"template1", "template2"},
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    mock_adaptors.NewMockEnv(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-o", "hello", "-t", "template1 template2"}, version)
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})
}
