// +build !linux,!darwin,!windows

package executors

import (
	"errors"
	"os"
)

var (
	Supported = false
	Uid       = 0
	Gid       = 0
	SIGUSR1   = os.Interrupt
	SIGUSR2   = os.Interrupt
	SIGTERM   = os.Kill
)

func Move(dst, src string) error {
	return errors.New("Not supported")
}

func Chmod(f *os.File, perms os.FileMode) error {
	return errors.New("Not supported")
}

func Chown(f *os.File, uid, gid int) error {
	return errors.New("Not supported")
}
