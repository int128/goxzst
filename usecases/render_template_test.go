package usecases

import (
	"bytes"
	"io"
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/pkg/errors"
)

func TestRenderTemplate_Do(t *testing.T) {
	t.Run("homebrew.rb", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			LookupEnv("VERSION").
			Return("v1.0.0", true).
			AnyTimes()

		var b mock_adaptors.WriteBuffer
		filesystem := mock_adaptors.NewMockFileSystem(ctrl)
		filesystem.EXPECT().
			MkdirAll("dist")
		filesystem.EXPECT().
			Create("dist/output").
			Return(&b, nil)
		filesystem.EXPECT().
			Open("dist/output_darwin_amd64.zip").
			Return(ioutil.NopCloser(strings.NewReader("text1")), nil)

		u := RenderTemplate{
			Env:        env,
			FileSystem: filesystem,
			Logger:     mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.RenderTemplateIn{
			InputFilename:  "testdata/homebrew.rb",
			OutputFilename: "dist/output",
			Variables: map[string]string{
				"darwin_amd64_archive": "dist/output_darwin_amd64.zip",
			},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}

		want, err := ioutil.ReadFile("testdata/homebrew.want.rb")
		if err != nil {
			t.Fatalf("could not read want: %s", err)
		}
		if !bytes.Equal(want, b.Bytes()) {
			t.Errorf("rendered content wants \n----\n%s\n----\nbut \n----\n%s", string(want), b.String())
		}
	})

	t.Run("krew.yaml", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			LookupEnv("VERSION").
			Return("v1.0.0", true).
			AnyTimes()

		var b mock_adaptors.WriteBuffer
		filesystem := mock_adaptors.NewMockFileSystem(ctrl)
		filesystem.EXPECT().
			MkdirAll("dist")
		filesystem.EXPECT().
			Create("dist/output").
			Return(&b, nil)
		filesystem.EXPECT().
			Open(gomock.Any()).
			DoAndReturn(func(name string) (io.ReadCloser, error) {
				switch name {
				case "dist/output_linux_amd64.zip":
					return ioutil.NopCloser(strings.NewReader("text1")), nil
				case "dist/output_darwin_amd64.zip":
					return ioutil.NopCloser(strings.NewReader("text2")), nil
				case "dist/output_windows_amd64.zip":
					return ioutil.NopCloser(strings.NewReader("text3")), nil
				}
				return nil, errors.Errorf("no such file: %s", name)
			}).
			AnyTimes()

		u := RenderTemplate{
			Env:        env,
			FileSystem: filesystem,
			Logger:     mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.RenderTemplateIn{
			InputFilename:  "testdata/krew.yaml",
			OutputFilename: "dist/output",
			Variables: map[string]string{
				"linux_amd64_archive":   "dist/output_linux_amd64.zip",
				"darwin_amd64_archive":  "dist/output_darwin_amd64.zip",
				"windows_amd64_archive": "dist/output_windows_amd64.zip",
			},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}

		want, err := ioutil.ReadFile("testdata/krew.want.yaml")
		if err != nil {
			t.Fatalf("could not read want: %s", err)
		}
		if !bytes.Equal(want, b.Bytes()) {
			t.Errorf("rendered content wants \n----\n%s\n----\nbut \n----\n%s", string(want), b.String())
		}
	})

	t.Run("NoSuchEnv", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		env := mock_adaptors.NewMockEnv(ctrl)
		env.EXPECT().
			LookupEnv("VERSION").
			Return("", false).
			AnyTimes()

		var b mock_adaptors.WriteBuffer
		filesystem := mock_adaptors.NewMockFileSystem(ctrl)
		filesystem.EXPECT().
			MkdirAll("dist")
		filesystem.EXPECT().
			Create("dist/output").
			Return(&b, nil)

		u := RenderTemplate{
			Env:        env,
			FileSystem: filesystem,
			Logger:     mock_adaptors.NewLogger(t),
		}
		err := u.Do(usecases.RenderTemplateIn{
			InputFilename:  "testdata/homebrew.rb",
			OutputFilename: "dist/output",
			Variables: map[string]string{
				"darwin_amd64_archive": "dist/output_linux_amd64.zip",
			},
		})
		if err == nil {
			t.Errorf("Do wants error but none")
		}
		t.Logf("Do returned expected error: %s", err)
	})
}
