package archiver

import (
	"errors"
	"path/filepath"
	"regexp"
)

var archivedRegex = regexp.MustCompile(`(\.tar|\.tar\.gz|\.zip)$`)

func Extension(path string) string {
	ext := archivedRegex.FindStringSubmatch(path)
	if len(ext) == 0 {
		return ""
	}
	return ext[1]
}

func CheckPath(path string) error {
	if len(path) == 0 {
		return errors.New("empty path")
	} else if filepath.IsAbs(path) {
		return errors.New("non-relative path: " + path)
	}
	return nil
}
