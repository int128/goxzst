package usecases

import (
	"archive/zip"
	"bytes"
	"io"
	"io/ioutil"
	"os"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/mock_adaptors"
	"github.com/int128/goxzst/usecases/interfaces"
)

func TestArchive_Do(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	var zipBuffer bytes.Buffer
	filesystem := mock_adaptors.NewMockFilesystem(ctrl)
	filesystem.EXPECT().
		Create("output").
		Return(&nopWriteCloser{&zipBuffer}, nil)
	filesystem.EXPECT().
		GetMode("input1").
		Return(os.FileMode(0644), nil)
	filesystem.EXPECT().
		Open("input1").
		Return(ioutil.NopCloser(strings.NewReader("text1")), nil)
	filesystem.EXPECT().
		GetMode("input2").
		Return(os.FileMode(0755), nil)
	filesystem.EXPECT().
		Open("input2").
		Return(ioutil.NopCloser(strings.NewReader("text2")), nil)

	u := Archive{
		Filesystem: filesystem,
	}
	if err := u.Do(usecases.ArchiveIn{
		OutputFilename: "output",
		Entries: []usecases.ArchiveEntry{
			{
				Path:          "entry1",
				InputFilename: "input1",
			}, {
				Path:          "entry2",
				InputFilename: "input2",
			},
		},
	}); err != nil {
		t.Errorf("Do returned error: %+v", err)
	}

	r, err := zip.NewReader(bytes.NewReader(zipBuffer.Bytes()), int64(zipBuffer.Len()))
	if err != nil {
		t.Fatalf("error while reading created zip: %s", err)
	}
	if len(r.File) != 2 {
		t.Errorf("len wants 2 but %d", len(r.File))
	}
	assertZipEntry(t, r.File[0], "entry1", 0644, []byte("text1"))
	assertZipEntry(t, r.File[1], "entry2", 0755, []byte("text2"))
}

func assertZipEntry(t *testing.T, f *zip.File, name string, mode os.FileMode, content []byte) {
	if f.Name != name {
		t.Errorf("Name wants %s but %s", name, f.Name)
	}
	if f.Mode() != mode {
		t.Errorf("Mode wants %v but %v", mode, f.Mode())
	}
	r, err := f.Open()
	if err != nil {
		t.Errorf("error while opening zip entry: %s", err)
		return
	}
	defer r.Close()
	c, err := ioutil.ReadAll(r)
	if err != nil {
		t.Errorf("error while reading zip entry: %s", err)
		return
	}
	if bytes.Compare(c, content) != 0 {
		t.Errorf("content wants %v but %v", content, c)
	}
}

type nopWriteCloser struct {
	io.Writer
}

func (*nopWriteCloser) Close() error {
	return nil
}
