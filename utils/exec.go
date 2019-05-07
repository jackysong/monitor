package utils

import "os/exec"

func ExecShell(command string) (string, error) {
	cmd := exec.Command("/bin/sh", "-c", command)

	out, err := cmd.CombinedOutput()
	if ErrorCheck(err) {
		return "", err
	}
	return string(out), err
}

func ExecExpect(command string) (string, error) {
	cmd := exec.Command("/usr/bin/expect", "-c", command)

	out, err := cmd.Output()
	if ErrorCheck(err) {
		return "", err
	}
	return string(out), err
}
