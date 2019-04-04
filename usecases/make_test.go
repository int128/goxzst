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
				Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
			})
		archive := mock_usecases.NewMockArchive(ctrl)
		archive.EXPECT().
			Do(usecases.ArchiveIn{
				OutputFilename: "output_linux_amd64.zip",
				Entries: []usecases.ArchiveEntry{
					{Path: "output", InputFilename: "output_linux_amd64"},
				},
			})
		digest := mock_usecases.NewMockDigest(ctrl)
		digest.EXPECT().
			Do(usecases.DigestIn{
				OutputFilename: "output_linux_amd64.zip.sha256",
				InputFilename:  "output_linux_amd64.zip",
			}).
			Return(&usecases.DigestOut{SHA256: "sha256"}, nil)

		filesystem := mock_adaptors.NewMockFilesystem(ctrl)
		filesystem.EXPECT().
			Remove("output_linux_amd64")

		u := Make{
			CrossBuild:     crossBuild,
			Archive:        archive,
			Digest:         digest,
			RenderTemplate: mock_usecases.NewMockRenderTemplate(ctrl),
			Filesystem:     filesystem,
			Logger:         mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.MakeIn{
			OutputName: "output",
			Platforms:  []build.Platform{{GOOS: "linux", GOARCH: "amd64"}},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("FullOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		filesystem := mock_adaptors.NewMockFilesystem(ctrl)
		filesystem.EXPECT().
			Remove("dir/output_linux_amd64")

		crossBuild := mock_usecases.NewMockCrossBuild(ctrl)
		crossBuild.EXPECT().
			Do(usecases.CrossBuildIn{
				OutputFilename: "dir/output_linux_amd64",
				Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		archive := mock_usecases.NewMockArchive(ctrl)
		archive.EXPECT().
			Do(usecases.ArchiveIn{
				OutputFilename: "dir/output_linux_amd64.zip",
				Entries: []usecases.ArchiveEntry{
					{Path: "output", InputFilename: "dir/output_linux_amd64"},
					{Path: "LICENSE", InputFilename: "LICENSE"},
				},
			})
		digest := mock_usecases.NewMockDigest(ctrl)
		digest.EXPECT().
			Do(usecases.DigestIn{
				OutputFilename: "dir/output_linux_amd64.zip.sha256",
				InputFilename:  "dir/output_linux_amd64.zip",
			}).
			Return(&usecases.DigestOut{SHA256: "sha256"}, nil)
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
			Archive:        archive,
			Digest:         digest,
			RenderTemplate: renderTemplate,
			Filesystem:     filesystem,
			Logger:         mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.MakeIn{
			OutputDir:  "dir",
			OutputName: "output",
			Platforms: []build.Platform{
				{GOOS: "linux", GOARCH: "amd64"},
			},
			GoBuildArgs:           []string{"-ldflags", "-X foo=bar"},
			ArchiveExtraFilenames: []string{"LICENSE"},
			TemplateFilenames:     []string{"template1"},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})
}
