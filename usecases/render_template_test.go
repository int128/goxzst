package usecases

import (
	"bytes"
	"io/ioutil"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/usecases/interfaces"
)

func TestNewRenderTemplate(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var b bytes.Buffer
	filesystem := mock_adaptors.NewMockFilesystem(ctrl)
	filesystem.EXPECT().
		Create("output").
		Return(&nopWriteCloser{&b}, nil)

	u := RenderTemplate{
		Filesystem: filesystem,
	}
	if err := u.Do(usecases.RenderTemplateIn{
		InputFilename:  "testdata/goxzst.rb",
		OutputFilename: "output",
		Variables: map[string]string{
			"version":                 "v1.0.0",
			"darwin_amd64_zip_sha256": "0123456789abcdef0123456789abcdef0123456789abcdef0123456789abcdef",
		},
	}); err != nil {
		t.Errorf("Do returned error: %+v", err)
	}

	want, err := ioutil.ReadFile("testdata/goxzst.want.rb")
	if err != nil {
		t.Fatalf("could not read want: %s", err)
	}
	if bytes.Compare(want, b.Bytes()) != 0 {
		t.Errorf("rendered content wants \n----\n%s\n----\nbut \n----\n%s", string(want), b.String())
	}
}
