package xzs

import (
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/google/go-cmp/cmp"
	"github.com/int128/goxzst/models/build"
	"github.com/int128/goxzst/models/digest"
	"github.com/int128/goxzst/usecases/archive"
	"github.com/int128/goxzst/usecases/archive/mock_archive"
	"github.com/int128/goxzst/usecases/crossbuild"
	"github.com/int128/goxzst/usecases/crossbuild/mock_crossbuild"
	digestUseCase "github.com/int128/goxzst/usecases/digest"
	"github.com/int128/goxzst/usecases/digest/mock_digest"
)

func TestMakeSingle_Do(t *testing.T) {
	t.Run("LessOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		platform := build.Platform{GOOS: "linux", GOARCH: "amd64"}
		mockCrossBuild := mock_crossbuild.NewMockInterface(ctrl)
		mockCrossBuild.EXPECT().
			Do(crossbuild.Input{
				OutputFilename: "output_linux_amd64",
				Platform:       platform,
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

		u := XZS{
			CrossBuild: mockCrossBuild,
			Archive:    mockArchive,
			Digest:     mockDigest,
		}
		artifact, err := u.Do(Input{
			OutputName:      "output",
			Platform:        platform,
			DigestAlgorithm: digest.SHA256,
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
		want := &build.Artifact{
			Platform:       platform,
			ExecutableFile: build.ExecutableFile{Base: "output_linux_amd64", Platform: platform},
			ArchiveFile:    build.ArchiveFile{Base: "output_linux_amd64", Suffix: ".zip"},
			DigestFile:     build.DigestFile{Base: "output_linux_amd64.zip", Suffix: ".sha256"},
		}
		if diff := cmp.Diff(want, artifact); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})

	t.Run("FullOptions", func(t *testing.T) {
		ctrl := gomock.NewController(t)
		defer ctrl.Finish()

		platform := build.Platform{GOOS: "windows", GOARCH: "amd64"}
		mockCrossBuild := mock_crossbuild.NewMockInterface(ctrl)
		mockCrossBuild.EXPECT().
			Do(crossbuild.Input{
				OutputFilename: "dir/output_windows_amd64.exe",
				Platform:       platform,
				GoBuildArgs:    []string{"-ldflags", "-X foo=bar"},
			})
		mockArchive := mock_archive.NewMockInterface(ctrl)
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
				OutputFilename: "dir/output_windows_amd64.zip.sha256",
				InputFilename:  "dir/output_windows_amd64.zip",
				Algorithm:      digest.SHA256,
			})

		u := XZS{
			CrossBuild: mockCrossBuild,
			Archive:    mockArchive,
			Digest:     mockDigest,
		}
		artifact, err := u.Do(Input{
			OutputDir:             "dir",
			OutputName:            "output",
			Platform:              platform,
			GoBuildArgs:           []string{"-ldflags", "-X foo=bar"},
			ArchiveExtraFilenames: []string{"LICENSE"},
			DigestAlgorithm:       digest.SHA256,
		})
		if err != nil {
			t.Errorf("Do returned error: %+v", err)
		}
		want := &build.Artifact{
			Platform:       platform,
			ExecutableFile: build.ExecutableFile{Base: "dir/output_windows_amd64", Platform: platform},
			ArchiveFile:    build.ArchiveFile{Base: "dir/output_windows_amd64", Suffix: ".zip"},
			DigestFile:     build.DigestFile{Base: "dir/output_windows_amd64.zip", Suffix: ".sha256"},
		}
		if diff := cmp.Diff(want, artifact); diff != "" {
			t.Errorf("mismatch (-want +got):\n%s", diff)
		}
	})
}
