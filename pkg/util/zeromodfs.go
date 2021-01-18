package util

import (
	"net/http"
	"os"
	"time"
)

// ZeroModTimeFileSystem is an http.FileSystem wrapper.
// It exposes a filesystem exactly like Source, except
// all file modification times are changed to zero.
type ZeroModTimeFileSystem struct {
	Source http.FileSystem
}

func (fs ZeroModTimeFileSystem) Open(name string) (http.File, error) {
	f, err := fs.Source.Open(name)

	return file{f}, err
}

type file struct {
	http.File
}

func (f file) Stat() (os.FileInfo, error) {
	fi, err := f.File.Stat()

	return fileInfo{fi}, err
}

type fileInfo struct {
	os.FileInfo
}

func (fi fileInfo) ModTime() time.Time { return time.Time{} }
