// Package build contains the models of build.
package build

import "fmt"

type GOOS string
type GOARCH string

// Platform represents a platform to run an executable.
type Platform struct {
	GOOS   GOOS
	GOARCH GOARCH
}

func (p Platform) String() string {
	return fmt.Sprintf("%s_%s", p.GOOS, p.GOARCH)
}

// ExecutableFile represents a file that is an outcome of go build.
type ExecutableFile struct {
	Base     string
	Platform Platform
}

// Suffix returns the platform dependent suffix.
func (f *ExecutableFile) Suffix() string {
	if f.Platform.GOOS == "windows" {
		return ".exe"
	}
	return ""
}

// Name returns the basename and the suffix for the platform.
func (f *ExecutableFile) Name() string {
	return f.Base + f.Suffix()
}

// ArchiveFile represents a file which has an executable file and extra files.
type ArchiveFile struct {
	Base   string
	Suffix string
}

// Name returns the basename and the suffix.
func (f *ArchiveFile) Name() string {
	return f.Base + f.Suffix
}

// DigestFile represents a file which has a digest of an archive file.
type DigestFile struct {
	Base   string
	Suffix string
}

// Name returns the basename and the suffix.
func (f *DigestFile) Name() string {
	return f.Base + f.Suffix
}

// Artifact represents a set of an executable, archive and digest file.
type Artifact struct {
	Platform       Platform
	ExecutableFile ExecutableFile
	ArchiveFile    ArchiveFile
	DigestFile     DigestFile
}
