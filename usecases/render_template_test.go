package usecases

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/usecases/interfaces"
)

func TestNewRenderTemplate(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
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

		u := RenderTemplate{
			Env:        env,
			FileSystem: filesystem,
			Logger:     mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.RenderTemplateIn{
			InputFilename:  "testdata/homebrew.rb",
			OutputFilename: "dist/output",
			Variables: map[string]string{
				"darwin_amd64_zip_sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
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
				"darwin_amd64_zip_sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
			},
		})
		if err == nil {
			t.Errorf("Do wants error but none")
		}
		t.Logf("Do returned expected error: %s", err)
	})
}
