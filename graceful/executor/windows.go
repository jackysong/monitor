// +build windows

package executor

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"syscall"
)

var (
	Supported = true
	Uid       = syscall.Getuid()
	Gid       = syscall.Getgid()
	SIGUSR1   = syscall.SIGTERM
	SIGUSR2   = syscall.SIGTERM
	SIGTERM   = syscall.SIGTERM
)

func Move(dst, src string) error {
	os.MkdirAll(filepath.Dir(dst), 0755)
	if err := os.Rename(src, dst); err == nil {
		return nil
	}
	//HACK: we're shelling out to move because windows
	//throws errors when crossing device boundaryes.
	// https://www.microsoft.com/resources/documentation/windows/xp/all/proddocs/en-us/move.mspx?mfr=true

	// https://blogs.msdn.microsoft.com/twistylittlepassagesallalike/2011/04/23/everyone-quotes-command-line-arguments-the-wrong-way/
	R := func(s string) string { return replShellMeta.Replace(syscall.EscapeArg(s)) }
	cmd := exec.Command("cmd", "/c", `move /y `+R(src)+` `+R(dst))
	if b, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%v: %q: %v", cmd.Args, bytes.TrimSpace(b), err)
	}
	return nil
}

func Chmod(f *os.File, perms os.FileMode) error {
	if err := f.Chmod(perms); err != nil && !strings.Contains(err.Error(), "not supported") {
		return err
	}
	return nil
}

func Chown(f *os.File, uid, gid int) error {
	if err := f.Chown(uid, gid); err != nil && !strings.Contains(err.Error(), "not supported") {
		return err
	}
	return nil
}

// https://blogs.msdn.microsoft.com/twistylittlepassagesallalike/2011/04/23/everyone-quotes-command-line-arguments-the-wrong-way/
var replShellMeta = strings.NewReplacer(
	`(`, `^(`,
	`)`, `^)`,
	`%`, `^%`,
	`!`, `^!`,
	`^`, `^^`,
	`"`, `^"`,
	`<`, `^<`,
	`>`, `^>`,
	`&`, `^&`,
	`|`, `^|`,
)
