package digest

import (
	"io/ioutil"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/fs/mock_fs"
	testingFs "github.com/int128/goxzst/adaptors/fs/testing"
	testingLogger "github.com/int128/goxzst/adaptors/logger/testing"
	"github.com/int128/goxzst/models/digest"
)

func TestDigest_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var b testingFs.WriteBuffer
	filesystem := mock_fs.NewMockInterface(ctrl)
	filesystem.EXPECT().
		MkdirAll("dist")
	filesystem.EXPECT().
		Create("dist/output").
		Return(&b, nil)
	filesystem.EXPECT().
		Open("input1").
		Return(ioutil.NopCloser(strings.NewReader("text1")), nil)

	u := Digest{
		FileSystem: filesystem,
		Logger:     testingLogger.New(t),
	}
	if err := u.Do(Input{
		InputFilename:  "input1",
		OutputFilename: "dist/output",
		Algorithm:      digest.SHA256,
	}); err != nil {
		t.Errorf("Do returned error: %+v", err)
	}

	// echo -n text1 | shasum -a 256
	const text1SHA256 = "fe8df1a5a1980493ca9406ad3bb0e41297d979d90165a181fb39a5616a1c0789"
	if b.String() != text1SHA256 {
		t.Errorf("output file content wants %s but %s", text1SHA256, b.String())
	}
}
