package usecases

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/usecases/interfaces"
	"github.com/int128/goxzst/usecases/mock_usecases"
)

func TestMake_Do(t *testing.T) {
	t.Run("LessOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		crossBuild := mock_usecases.NewMockCrossBuild(ctrl)
		crossBuild.EXPECT().
			Do(usecases.CrossBuildIn{
				OutputFilename: "output_linux_amd64",
				Target:         build.Target{GOOS: "linux", GOARCH: "amd64"},
			})
		createZip := mock_usecases.NewMockCreateZip(ctrl)
		createZip.EXPECT().
			Do(usecases.CreateZipIn{
				OutputFilename: "output_linux_amd64.zip",
				Entries: []usecases.ZipEntry{
					{Path: "output", InputFilename: "output_linux_amd64"},
				},
			})
		createSHA := mock_usecases.NewMockCreateSHA(ctrl)
		createSHA.EXPECT().
			Do(usecases.CreateSHAIn{
				OutputFilename: "output_linux_amd64.zip.sha256",
				InputFilename:  "output_linux_amd64.zip",
			}).
			Return(&usecases.CreateSHAOut{SHA256: "sha256"}, nil)

		u := Make{
			CrossBuild:     crossBuild,
			CreateZip:      createZip,
			CreateSHA:      createSHA,
			RenderTemplate: mock_usecases.NewMockRenderTemplate(ctrl),
			Filesystem:     mock_adaptors.NewMockFilesystem(ctrl),
		}
		if err := u.Do(usecases.MakeIn{
			OutputName: "output",
			Targets:    []build.Target{{GOOS: "linux", GOARCH: "amd64"}},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("FullOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filesystem := mock_adaptors.NewMockFilesystem(ctrl)
		filesystem.EXPECT().
			MkdirAll("dir")

		crossBuild := mock_usecases.NewMockCrossBuild(ctrl)
		crossBuild.EXPECT().
			Do(usecases.CrossBuildIn{
				OutputFilename: "dir/output_linux_amd64",
				Target:         build.Target{GOOS: "linux", GOARCH: "amd64"},
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		createZip := mock_usecases.NewMockCreateZip(ctrl)
		createZip.EXPECT().
			Do(usecases.CreateZipIn{
				OutputFilename: "dir/output_linux_amd64.zip",
				Entries: []usecases.ZipEntry{
					{Path: "output", InputFilename: "dir/output_linux_amd64"},
				},
			})
		createSHA := mock_usecases.NewMockCreateSHA(ctrl)
		createSHA.EXPECT().
			Do(usecases.CreateSHAIn{
				OutputFilename: "dir/output_linux_amd64.zip.sha256",
				InputFilename:  "dir/output_linux_amd64.zip",
			}).
			Return(&usecases.CreateSHAOut{SHA256: "sha256"}, nil)
		renderTemplate := mock_usecases.NewMockRenderTemplate(ctrl)
		renderTemplate.EXPECT().
			Do(usecases.RenderTemplateIn{
				InputFilename:  "template1",
				OutputFilename: "dir/template1",
				Variables: map[string]string{
					"linux_amd64_zip_sha256": "sha256",
				},
			})

		u := Make{
			CrossBuild:     crossBuild,
			CreateZip:      createZip,
			CreateSHA:      createSHA,
			RenderTemplate: renderTemplate,
			Filesystem:     filesystem,
		}
		if err := u.Do(usecases.MakeIn{
			OutputDir:  "dir",
			OutputName: "output",
			Targets: []build.Target{
				{GOOS: "linux", GOARCH: "amd64"},
			},
			GoBuildArgs:       []string{"-ldflags", "-X foo=bar"},
			TemplateFilenames: []string{"template1"},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})
}
