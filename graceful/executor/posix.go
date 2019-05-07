// +build linux darwin

package executor

//this file attempts to contain all posix
//specific stuff, that needs to be implemented
//in some other way on other OSs... TODO!

import (
	"os"
	"os/exec"
	"syscall"
)

var (
	Supported = true
	Uid       = syscall.Getuid()
	Gid       = syscall.Getgid()
	SIGUSR1   = syscall.SIGUSR1
	SIGUSR2   = syscall.SIGUSR2
	SIGTERM   = syscall.SIGTERM
)

func Move(dst, src string) error {
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	//HACK: we're shelling out to mv because linux
	//throws errors when crossing device boundaryes.
	//TODO see sys_posix_mv.go
	return exec.Command("mv", src, dst).Run()
}

func Chmod(f *os.File, perms os.FileMode) error {
	return f.Chmod(perms)
}
func Chown(f *os.File, uid, gid int) error {
	return f.Chown(uid, gid)
}
