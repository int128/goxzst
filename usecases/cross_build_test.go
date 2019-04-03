package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/interfaces"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/usecases/interfaces"
)

func TestCrossBuild_Do(t *testing.T) {
	t.Run("BasicOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()
		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Exec(adaptors.ExecIn{
				Name:     "go",
				Args:     []string{"build", "-o", "output"},
				ExtraEnv: []string{"GOOS=linux", "GOARCH=amd64"},
			})

		u := CrossBuild{
			Env:    env,
			Logger: mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.CrossBuildIn{
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
		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			Exec(adaptors.ExecIn{
				Name:     "go",
				Args:     []string{"build", "-o", "output", "-ldflags", "-X foo=bar"},
				ExtraEnv: []string{"GOOS=linux", "GOARCH=amd64"},
			})

		u := CrossBuild{
			Env:    env,
			Logger: mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.CrossBuildIn{
			OutputFilename: "output",
			GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
		}); err != nil {
			t.Fatalf("Do returned error: %+v", err)
		}
	})
}
