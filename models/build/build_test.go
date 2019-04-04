package build_test

import (
	"testing"

	"github.com/int128/goxzst/models/build"
)

func TestExecutableFile_Name(t *testing.T) {
	t.Run("non-windows", func(t *testing.T) {
		f := build.ExecutableFile{
			Base:     "hello",
			Platform: build.Platform{GOOS: "linux", GOARCH: "amd64"},
		}
		name := f.Name()
		want := "hello"
		if name != want {
			t.Errorf("Name wants %s but %s", want, name)
		}
	})

	t.Run("windows", func(t *testing.T) {
		f := build.ExecutableFile{
			Base:     "hello",
			Platform: build.Platform{GOOS: "windows", GOARCH: "amd64"},
		}
		name := f.Name()
		want := "hello.exe"
		if name != want {
			t.Errorf("Name wants %s but %s", want, name)
		}
	})
}

func TestArchiveFile_Name(t *testing.T) {
	f := build.ArchiveFile{
		Base:   "hello",
		Suffix: ".zip",
	}
	name := f.Name()
	want := "hello.zip"
	if name != want {
		t.Errorf("Name wants %s but %s", want, name)
	}
}

func TestDigestFile_Name(t *testing.T) {
	f := build.DigestFile{
		Base:   "hello",
		Suffix: ".sha256",
	}
	name := f.Name()
	want := "hello.sha256"
	if name != want {
		t.Errorf("Name wants %s but %s", want, name)
	}
}
