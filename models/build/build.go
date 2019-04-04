package build

import "fmt"

type GOOS string
type GOARCH string

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

type ArchiveFile struct {
	Base   string
	Suffix string
}

// Name returns the basename and the suffix.
func (f *ArchiveFile) Name() string {
	return f.Base + f.Suffix
}

type DigestFile struct {
	Base   string
	Suffix string
}

// Name returns the basename and the suffix.
func (f *DigestFile) Name() string {
	return f.Base + f.Suffix
}
