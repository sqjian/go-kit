package s3

import (
	"io"
	"io/fs"
	"time"
)

// Fake implementation of files.

func NewFsFile(name string, contents []byte, mode fs.FileMode) fs.File {
	return &FakeFile{
		name:     name,
		contents: contents,
		mode:     mode,
	}
}

// FakeFile implements FileLike and also fs.FileInfo.
type FakeFile struct {
	name     string
	contents []byte
	mode     fs.FileMode
	offset   int
}

// Reset prepares a FakeFile for reuse.
func (f *FakeFile) Reset() *FakeFile {
	f.offset = 0
	return f
}

// FileLike methods.

func (f *FakeFile) Name() string {
	// A bit of a cheat: we only have a basename, so that's also ok for FileInfo.
	return f.name
}

func (f *FakeFile) Stat() (fs.FileInfo, error) {
	return f, nil
}

func (f *FakeFile) Read(p []byte) (int, error) {
	if f.offset >= len(f.contents) {
		return 0, io.EOF
	}
	n := copy(p, f.contents[f.offset:])
	f.offset += n
	return n, nil
}

func (f *FakeFile) Close() error {
	return nil
}

// fs.FileInfo methods.

func (f *FakeFile) Size() int64 {
	return int64(len(f.contents))
}

func (f *FakeFile) Mode() fs.FileMode {
	return f.mode
}

func (f *FakeFile) ModTime() time.Time {
	return time.Time{}
}

func (f *FakeFile) IsDir() bool {
	return false
}

func (f *FakeFile) Sys() any {
	return nil
}
