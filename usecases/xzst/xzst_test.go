package xzst

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/fs/mock_fs"
	testingLogger "github.com/int128/goxzst/adaptors/logger/testing"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/rendertemplate"
	"github.com/int128/goxzst/usecases/rendertemplate/mock_rendertemplate"
	"github.com/int128/goxzst/usecases/xzs"
	"github.com/int128/goxzst/usecases/xzs/mock_xzs"
)

func TestMake_Do(t *testing.T) {
	t.Run("LessOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		platform := build.Platform{GOOS: "linux", GOARCH: "amd64"}
		xzsUseCase := mock_xzs.NewMockInterface(ctrl)
		xzsUseCase.EXPECT().
			Do(xzs.Input{
				OutputName:      "output",
				Platform:        platform,
				DigestAlgorithm: digest.SHA256,
			}).
			Return(&build.Artifact{
				Platform:       platform,
				ExecutableFile: build.ExecutableFile{Base: "output_linux_amd64", Platform: platform},
				ArchiveFile:    build.ArchiveFile{Base: "output_linux_amd64", Suffix: ".zip"},
				DigestFile:     build.DigestFile{Base: "output_linux_amd64.zip", Suffix: ".sha256"},
			}, nil)
		mockFileSystem := mock_fs.NewMockInterface(ctrl)
		mockFileSystem.EXPECT().
			Remove("output_linux_amd64")

		u := XZST{
			XZS:            xzsUseCase,
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

		linuxPlatform := build.Platform{GOOS: "linux", GOARCH: "amd64"}
		windowsPlatform := build.Platform{GOOS: "windows", GOARCH: "amd64"}
		xzsUseCase := mock_xzs.NewMockInterface(ctrl)
		xzsUseCase.EXPECT().
			Do(xzs.Input{
				OutputDir:             "dir",
				OutputName:            "output",
				Platform:              linuxPlatform,
				GoBuildArgs:           []string{"-ldflags", "-X foo=bar"},
				ArchiveExtraFilenames: []string{"LICENSE"},
				DigestAlgorithm:       digest.SHA256,
			}).
			Return(&build.Artifact{
				Platform:       linuxPlatform,
				ExecutableFile: build.ExecutableFile{Base: "dir/output_linux_amd64", Platform: linuxPlatform},
				ArchiveFile:    build.ArchiveFile{Base: "dir/output_linux_amd64", Suffix: ".zip"},
				DigestFile:     build.DigestFile{Base: "dir/output_linux_amd64.zip", Suffix: ".sha256"},
			}, nil)
		xzsUseCase.EXPECT().
			Do(xzs.Input{
				OutputDir:             "dir",
				OutputName:            "output",
				Platform:              windowsPlatform,
				GoBuildArgs:           []string{"-ldflags", "-X foo=bar"},
				ArchiveExtraFilenames: []string{"LICENSE"},
				DigestAlgorithm:       digest.SHA256,
			}).
			Return(&build.Artifact{
				Platform:       windowsPlatform,
				ExecutableFile: build.ExecutableFile{Base: "dir/output_windows_amd64", Platform: windowsPlatform},
				ArchiveFile:    build.ArchiveFile{Base: "dir/output_windows_amd64", Suffix: ".zip"},
				DigestFile:     build.DigestFile{Base: "dir/output_windows_amd64.zip", Suffix: ".sha256"},
			}, nil)
		mockFileSystem := mock_fs.NewMockInterface(ctrl)
		mockFileSystem.EXPECT().
			Remove("dir/output_linux_amd64")
		mockFileSystem.EXPECT().
			Remove("dir/output_windows_amd64.exe")
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

		u := XZST{
			XZS:            xzsUseCase,
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
