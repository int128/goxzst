package archive

import (
	"archive/zip"
	"bytes"
	"io/ioutil"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/int128/goxzst/adaptors/fs/mock_fs"
	testingFs "github.com/int128/goxzst/adaptors/fs/testing"
	testingLogger "github.com/int128/goxzst/adaptors/logger/testing"
)

func TestArchive_Do(t *testing.T) {
	fileInfo1 := testingFs.FileInfo{
		ModeValue:    0644,
		ModTimeValue: time.Date(2019, 4, 1, 2, 3, 4, 0, time.UTC),
	}
	fileInfo2 := testingFs.FileInfo{
		ModeValue:    0755,
		ModTimeValue: time.Date(2019, 4, 9, 8, 7, 6, 0, time.UTC),
	}

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
		Stat("input1").
		Return(&fileInfo1, nil)
	filesystem.EXPECT().
		Open("input1").
		Return(ioutil.NopCloser(strings.NewReader("text1")), nil)
	filesystem.EXPECT().
		Stat("input2").
		Return(&fileInfo2, nil)
	filesystem.EXPECT().
		Open("input2").
		Return(ioutil.NopCloser(strings.NewReader("text2")), nil)

	u := Archive{
		FileSystem: filesystem,
		Logger:     testingLogger.New(t),
	}
	if err := u.Do(Input{
		OutputFilename: "dist/output",
		Entries: []Entry{
			{
				Filename:      "entry1",
				InputFilename: "input1",
			}, {
				Filename:      "entry2",
				InputFilename: "input2",
			},
		},
	}); err != nil {
		t.Errorf("Do returned error: %+v", err)
	}

	r, err := zip.NewReader(bytes.NewReader(b.Bytes()), int64(b.Len()))
	if err != nil {
		t.Fatalf("error while reading created zip: %s", err)
	}
	if len(r.File) != 2 {
		t.Errorf("len wants 2 but %d", len(r.File))
	}
	assertZipEntry(t, r.File[0], "entry1", []byte("text1"), &fileInfo1)
	assertZipEntry(t, r.File[1], "entry2", []byte("text2"), &fileInfo2)
}

func assertZipEntry(t *testing.T, f *zip.File, name string, content []byte, fileInfo os.FileInfo) {
	if f.Name != name {
		t.Errorf("Name wants %s but %s", name, f.Name)
	}
	if f.Mode() != fileInfo.Mode() {
		t.Errorf("Mode wants %v but %v", fileInfo.Mode(), f.Mode())
	}
	if f.Modified.UTC() != fileInfo.ModTime().UTC() {
		t.Errorf("ModTime wants %v but %v", fileInfo.ModTime(), f.Modified)
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
	if !bytes.Equal(c, content) {
		t.Errorf("content wants %v but %v", content, c)
	}
}
