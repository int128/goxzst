package crossbuild

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/env"
	"github.com/int128/goxzst/adaptors/env/mock_env"
	"github.com/int128/goxzst/adaptors/fs/mock_fs"
	testingLogger "github.com/int128/goxzst/adaptors/logger/testing"
	"github.com/int128/goxzst/models/build"
)

func TestCrossBuild_Do(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockEnv := mock_env.NewMockInterface(ctrl)
		mockEnv.EXPECT().
			Exec(env.Exec{
				Name:     "go",
				Args:     []string{"build", "-o", "output"},
				ExtraEnv: []string{"GOOS=linux", "GOARCH=amd64"},
			})
		mockFs := mock_fs.NewMockInterface(ctrl)
		mockFs.EXPECT().
			MkdirAll(".")

		u := CrossBuild{
			Env:        mockEnv,
			FileSystem: mockFs,
			Logger:     testingLogger.New(t),
		}
		if err := u.Do(Input{
			OutputFilename: "output",
			GoBuildArgs:    nil,
			Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
		}); err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
	})

	t.Run("WithGoBuildArgs", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		mockEnv := mock_env.NewMockInterface(ctrl)
		mockEnv.EXPECT().
			Exec(env.Exec{
				Name:     "go",
				Args:     []string{"build", "-o", "dist/output", "-ldflags", "-X foo=bar"},
				ExtraEnv: []string{"GOOS=linux", "GOARCH=amd64"},
			})
		mockFs := mock_fs.NewMockInterface(ctrl)
		mockFs.EXPECT().
			MkdirAll("dist")

		u := CrossBuild{
			Env:        mockEnv,
			FileSystem: mockFs,
			Logger:     testingLogger.New(t),
		}
		if err := u.Do(Input{
			OutputFilename: "dist/output",
			GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
		}); err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
	})
}
