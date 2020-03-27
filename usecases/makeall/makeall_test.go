package makeall

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/fs/mock_fs"
	testingLogger "github.com/int128/goxzst/adaptors/logger/testing"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/archive"
	"github.com/int128/goxzst/usecases/archive/mock_archive"
	"github.com/int128/goxzst/usecases/crossbuild"
	"github.com/int128/goxzst/usecases/crossbuild/mock_crossbuild"
	digestUseCase "github.com/int128/goxzst/usecases/digest"
	"github.com/int128/goxzst/usecases/digest/mock_digest"
	"github.com/int128/goxzst/usecases/rendertemplate"
	"github.com/int128/goxzst/usecases/rendertemplate/mock_rendertemplate"
)

func TestMake_Do(t *testing.T) {
	t.Run("LessOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		mockCrossBuild := mock_crossbuild.NewMockInterface(ctrl)
		mockCrossBuild.EXPECT().
			Do(crossbuild.Input{
				OutputFilename: "output_linux_amd64",
				Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
			})
		mockArchive := mock_archive.NewMockInterface(ctrl)
		mockArchive.EXPECT().
			Do(archive.Input{
				OutputFilename: "output_linux_amd64.zip",
				Entries: []archive.Entry{
					{Filename: "output", InputFilename: "output_linux_amd64"},
				},
			})
		mockDigest := mock_digest.NewMockInterface(ctrl)
		mockDigest.EXPECT().
			Do(digestUseCase.Input{
				OutputFilename: "output_linux_amd64.zip.sha256",
				InputFilename:  "output_linux_amd64.zip",
				Algorithm:      digest.SHA256,
			})

		mockFileSystem := mock_fs.NewMockInterface(ctrl)
		mockFileSystem.EXPECT().
			Remove("output_linux_amd64")

		u := MakeAll{
			CrossBuild:     mockCrossBuild,
			Archive:        mockArchive,
			Digest:         mockDigest,
			RenderTemplate: mock_rendertemplate.NewMockInterface(ctrl),
			FileSystem:     mockFileSystem,
			Logger:         testingLogger.New(t),
		}
		if err := u.Do(Input{
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

		mockFileSystem := mock_fs.NewMockInterface(ctrl)
		mockFileSystem.EXPECT().
			Remove("dir/output_linux_amd64")
		mockFileSystem.EXPECT().
			Remove("dir/output_windows_amd64.exe")

		mockCrossBuild := mock_crossbuild.NewMockInterface(ctrl)
		mockCrossBuild.EXPECT().
			Do(crossbuild.Input{
				OutputFilename: "dir/output_linux_amd64",
				Platform:       build.Platform{GOOS: "linux", GOARCH: "amd64"},
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		mockCrossBuild.EXPECT().
			Do(crossbuild.Input{
				OutputFilename: "dir/output_windows_amd64.exe",
				Platform:       build.Platform{GOOS: "windows", GOARCH: "amd64"},
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		mockArchive := mock_archive.NewMockInterface(ctrl)
		mockArchive.EXPECT().
			Do(archive.Input{
				OutputFilename: "dir/output_linux_amd64.zip",
				Entries: []archive.Entry{
					{Filename: "output", InputFilename: "dir/output_linux_amd64"},
					{Filename: "LICENSE", InputFilename: "LICENSE"},
				},
			})
		mockArchive.EXPECT().
			Do(archive.Input{
				OutputFilename: "dir/output_windows_amd64.zip",
				Entries: []archive.Entry{
					{Filename: "output.exe", InputFilename: "dir/output_windows_amd64.exe"},
					{Filename: "LICENSE", InputFilename: "LICENSE"},
				},
			})
		mockDigest := mock_digest.NewMockInterface(ctrl)
		mockDigest.EXPECT().
			Do(digestUseCase.Input{
				OutputFilename: "dir/output_linux_amd64.zip.sha256",
				InputFilename:  "dir/output_linux_amd64.zip",
				Algorithm:      digest.SHA256,
			})
		mockDigest.EXPECT().
			Do(digestUseCase.Input{
				OutputFilename: "dir/output_windows_amd64.zip.sha256",
				InputFilename:  "dir/output_windows_amd64.zip",
				Algorithm:      digest.SHA256,
			})
		mockRenderTemplate := mock_rendertemplate.NewMockInterface(ctrl)
		mockRenderTemplate.EXPECT().
			Do(rendertemplate.Input{
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

		u := MakeAll{
			CrossBuild:     mockCrossBuild,
			Archive:        mockArchive,
			Digest:         mockDigest,
			RenderTemplate: mockRenderTemplate,
			FileSystem:     mockFileSystem,
			Logger:         testingLogger.New(t),
		}
		if err := u.Do(Input{
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
