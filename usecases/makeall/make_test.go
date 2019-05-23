package makeall

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases"
	"github.com/int128/goxzst/usecases/mock_usecases"
)

func TestMake_Do(t *testing.T) {
	t.Run("LessOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrossBuild := mock_usecases.NewMockCrossBuild(ctrl)
		mockCrossBuild.EXPECT().
			Do(usecases.CrossBuildIn{
				OutputFilename: "output_linux_amd64",
				Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
			})
		mockArchive := mock_usecases.NewMockArchive(ctrl)
		mockArchive.EXPECT().
			Do(usecases.ArchiveIn{
				OutputFilename: "output_linux_amd64.zip",
				Entries: []usecases.ArchiveEntry{
					{Filename: "output", InputFilename: "output_linux_amd64"},
				},
			})
		mockDigest := mock_usecases.NewMockDigest(ctrl)
		mockDigest.EXPECT().
			Do(usecases.DigestIn{
				OutputFilename: "output_linux_amd64.zip.sha256",
				InputFilename:  "output_linux_amd64.zip",
				Algorithm:      digest.SHA256,
			})

		mockFileSystem := mock_adaptors.NewMockFileSystem(ctrl)
		mockFileSystem.EXPECT().
			Remove("output_linux_amd64")

		u := Make{
			CrossBuild:     mockCrossBuild,
			Archive:        mockArchive,
			Digest:         mockDigest,
			RenderTemplate: mock_usecases.NewMockRenderTemplate(ctrl),
			FileSystem:     mockFileSystem,
			Logger:         mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.MakeIn{
			OutputName:      "output",
			Platforms:       []build.Platform{{GOOS: "linux", GOARCH: "amd64"}},
			DigestAlgorithm: digest.SHA256,
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})

	t.Run("FullOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockFileSystem := mock_adaptors.NewMockFileSystem(ctrl)
		mockFileSystem.EXPECT().
			Remove("dir/output_linux_amd64")
		mockFileSystem.EXPECT().
			Remove("dir/output_windows_amd64.exe")

		mockCrossBuild := mock_usecases.NewMockCrossBuild(ctrl)
		mockCrossBuild.EXPECT().
			Do(usecases.CrossBuildIn{
				OutputFilename: "dir/output_linux_amd64",
				Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		mockCrossBuild.EXPECT().
			Do(usecases.CrossBuildIn{
				OutputFilename: "dir/output_windows_amd64.exe",
				Platform:       build.Platform{GOOS: "windows", GOARCH: "amd64"},
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		mockArchive := mock_usecases.NewMockArchive(ctrl)
		mockArchive.EXPECT().
			Do(usecases.ArchiveIn{
				OutputFilename: "dir/output_linux_amd64.zip",
				Entries: []usecases.ArchiveEntry{
					{Filename: "output", InputFilename: "dir/output_linux_amd64"},
					{Filename: "LICENSE", InputFilename: "LICENSE"},
				},
			})
		mockArchive.EXPECT().
			Do(usecases.ArchiveIn{
				OutputFilename: "dir/output_windows_amd64.zip",
				Entries: []usecases.ArchiveEntry{
					{Filename: "output.exe", InputFilename: "dir/output_windows_amd64.exe"},
					{Filename: "LICENSE", InputFilename: "LICENSE"},
				},
			})
		mockDigest := mock_usecases.NewMockDigest(ctrl)
		mockDigest.EXPECT().
			Do(usecases.DigestIn{
				OutputFilename: "dir/output_linux_amd64.zip.sha256",
				InputFilename:  "dir/output_linux_amd64.zip",
				Algorithm:      digest.SHA256,
			})
		mockDigest.EXPECT().
			Do(usecases.DigestIn{
				OutputFilename: "dir/output_windows_amd64.zip.sha256",
				InputFilename:  "dir/output_windows_amd64.zip",
				Algorithm:      digest.SHA256,
			})
		mockRenderTemplate := mock_usecases.NewMockRenderTemplate(ctrl)
		mockRenderTemplate.EXPECT().
			Do(usecases.RenderTemplateIn{
				InputFilename:  "template1",
				OutputFilename: "dir/template1",
				Variables: map[string]string{
					"linux_amd64_executable":   "dir/output_linux_amd64",
					"linux_amd64_archive":      "dir/output_linux_amd64.zip",
					"linux_amd64_digest":       "dir/output_linux_amd64.zip.sha256",
					"windows_amd64_executable": "dir/output_windows_amd64.exe",
					"windows_amd64_archive":    "dir/output_windows_amd64.zip",
					"windows_amd64_digest":     "dir/output_windows_amd64.zip.sha256",
				},
			})

		u := Make{
			CrossBuild:     mockCrossBuild,
			Archive:        mockArchive,
			Digest:         mockDigest,
			RenderTemplate: mockRenderTemplate,
			FileSystem:     mockFileSystem,
			Logger:         mock_adaptors.NewLogger(t),
		}
		if err := u.Do(usecases.MakeIn{
			OutputDir:  "dir",
			OutputName: "output",
			Platforms: []build.Platform{
				{GOOS: "linux", GOARCH: "amd64"},
				{GOOS: "windows", GOARCH: "amd64"},
			},
			GoBuildArgs:           []string{"-ldflags", "-X foo=bar"},
			ArchiveExtraFilenames: []string{"LICENSE"},
			DigestAlgorithm:       digest.SHA256,
			TemplateFilenames:     []string{"template1"},
		}); err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
	})
}
