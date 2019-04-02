package build

type GOOS string
type GOARCH string

type Target struct {
	GOOS   GOOS
	GOARCH GOARCH
}
