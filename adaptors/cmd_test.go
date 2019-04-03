package adaptors

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/int128/goxzst/usecases/mock_usecases"
)

func TestCmd_Run(t *testing.T) {
	defaultPlatforms := []build.Platform{
		{GOOS: "linux", GOARCH: "amd64"},
		{GOOS: "darwin", GOARCH: "amd64"},
		{GOOS: "windows", GOARCH: "amd64"},
	}

	t.Run("NoArgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		makeUseCase := mock_usecases.NewMockMake(ctrl)
		makeUseCase.EXPECT().
			Do(usecases.MakeIn{
				OutputDir:   "dist",
				OutputName:  "package",
				Platforms:   defaultPlatforms,
				GoBuildArgs: []string{},
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    newEnvMock(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst"})
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
				OutputDir:   "dist",
				OutputName:  "package",
				Platforms:   defaultPlatforms,
				GoBuildArgs: []string{"-ldflags", "-X foo=bar"},
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    newEnvMock(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "--", "-ldflags", "-X foo=bar"})
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
				OutputName: "package",
				Platforms: []build.Platform{
					{GOOS: "linux", GOARCH: "arm"},
				},
				GoBuildArgs: []string{},
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    newEnvMock(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-osarch", "linux_arm"})
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
				OutputName:        "package",
				Platforms:         defaultPlatforms,
				GoBuildArgs:       []string{},
				TemplateFilenames: []string{"template1", "template2"},
			})

		cmd := Cmd{
			Make:   makeUseCase,
			Logger: mock_adaptors.NewLogger(t),
			Env:    newEnvMock(ctrl),
		}
		exitCode := cmd.Run([]string{"goxzst", "-t", "template1 template2"})
		if exitCode != 0 {
			t.Errorf("exitCode wants 0 but %d", exitCode)
		}
	})
}

func newEnvMock(ctrl *gomock.Controller) *mock_adaptors.MockEnv {
	env := mock_adaptors.NewMockEnv(ctrl)
	env.EXPECT().
		Getwd().
		Return("/package", nil)
	return env
}
