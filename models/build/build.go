package build

type GOOS string
type GOARCH string

type Platform struct {
	GOOS   GOOS
	GOARCH GOARCH
}
