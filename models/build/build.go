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
