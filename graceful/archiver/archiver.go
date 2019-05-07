package archiver

import (
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"errors"
	"io"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var FileMode = os.FileMode(0755)

type action interface {
	addBytes(path string, contents []byte, mtime time.Time) error
	addFile(path string, info os.FileInfo, f *os.File) error
	addReader(path string, info os.FileInfo, r io.Reader) error
	close() error
}

type Archive struct {
	SizeLimit int64 //defaults to no limit (-1)
	FileLimit int   //defaults to no limit (-1)
	mutex     sync.Mutex
	writer    io.Writer
	action    action
}

//get writer with extension
func NewWriter(writer io.Writer, params ...interface{}) (*Archive, error) {
	archive := &Archive{
		SizeLimit: -1,
		FileLimit: -1,
		writer:    writer,
	}

	switch Extension(params[0].(string)) {
	case ".tar.gz":
		writer = gzip.NewWriter(writer)
		archive.writer = writer
		fallthrough
	case ".tar":
		archive.action = &tarActor{
			writer: tar.NewWriter(writer),
			mode:   FileMode,
		}
	case ".zip":
		var compressed bool
		if len(params) > 1 {
			compressed = params[1].(bool)
		}
		archive.action = &zipActor{
			writer:     zip.NewWriter(writer),
			compressed: compressed,
		}
	}

	return archive, nil
}

func (a *Archive) AddBytes(path string, contents []byte) error {
	return a.AddBytesMTime(path, contents, time.Now())
}

func (a *Archive) AddBytesMTime(path string, contents []byte, mtime time.Time) error {
	if err := CheckPath(path); err != nil {
		return err
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.action.addBytes(path, contents, mtime)
}

func (a *Archive) AddInfoReader(path string, info os.FileInfo, r io.Reader) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.action.addReader(path, info, r)
}

//You can prevent archive from performing an extra Stat by using AddInfoFile
//instead of AddFile
func (a *Archive) AddInfoFile(path string, info os.FileInfo, f *os.File) error {
	if err := CheckPath(path); err != nil {
		return err
	}
	a.mutex.Lock()
	defer a.mutex.Unlock()
	return a.action.addFile(path, info, f)
}

func (a *Archive) AddFile(path string, f *os.File) error {
	info, err := f.Stat()
	if err != nil {
		return err
	}
	return a.AddInfoFile(path, info, f)
}

func (a *Archive) AddDir(path string) error {
	size := int64(0)
	num := 0
	return filepath.Walk(path, func(p string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.Mode().IsRegular() {
			return nil
		}
		if a.SizeLimit >= 0 {
			size += info.Size()
			if size > a.SizeLimit {
				return errors.New("Surpassed maximum archive size")
			}
		}
		if a.FileLimit >= 0 {
			num++
			if num == a.FileLimit+1 {
				return errors.New("Surpassed maximum number of files in archive")
			}
		}
		rel, err := filepath.Rel(path, p)
		if err != nil {
			return err
		}
		f, err := os.Open(p)
		if err != nil {
			return err
		}
		defer f.Close()
		return a.action.addFile(rel, info, f)
	})
}

func (a *Archive) Close() error {
	if err := a.action.close(); err != nil {
		return err
	}
	if gz, ok := a.writer.(*gzip.Writer); ok {
		if err := gz.Close(); err != nil {
			return err
		}
	}
	return nil
}
