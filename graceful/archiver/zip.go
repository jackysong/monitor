package archiver

import (
	"archive/zip"
	"errors"
	"io"
	"os"
	"time"
)

type zipActor struct {
	writer     *zip.Writer
	compressed bool
}

func (a *zipActor) addReader(path string, info os.FileInfo, r io.Reader) error {
	if !info.Mode().IsRegular() {
		return errors.New("Only regular files supported: " + path)
	}

	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	h.Name = path
	if a.compressed {
		h.Method = zip.Deflate
	}

	zf, err := a.writer.CreateHeader(h)
	if err != nil {
		return err
	}

	_, err = io.Copy(zf, r)
	return err
}

func (a *zipActor) addBytes(path string, contents []byte, mtime time.Time) error {
	h := &zip.FileHeader{
		Name: path,
	}

	if a.compressed {
		h.Method = zip.Deflate
	}

	h.SetModTime(mtime)
	f, err := a.writer.CreateHeader(h)
	if err != nil {
		return err
	}
	_, err = f.Write(contents)
	return err
}

func (a *zipActor) addFile(path string, info os.FileInfo, f *os.File) error {
	if !info.Mode().IsRegular() {
		return errors.New("Only regular files supported: " + path)
	}
	h, err := zip.FileInfoHeader(info)
	if err != nil {
		return err
	}

	h.Name = path
	if a.compressed {
		h.Method = zip.Deflate
	}

	zf, err := a.writer.CreateHeader(h)
	if err != nil {
		return err
	}
	_, err = io.Copy(zf, f)
	return err
}

func (a *zipActor) close() error {
	return a.writer.Close()
}
