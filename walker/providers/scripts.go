package providers

import (
	"bytes"
	"os/exec"
	"strings"
)

func (s *Script) RunScript(address, flag string) (result string, err error) {
	output, _, err := runProcess(exec.Command("scripts/"+s.Name, address, flag))
	if err != nil {
		return "", err
	}
	return strings.TrimSuffix(output, "\n"), nil
}

func runProcess(process *exec.Cmd) (stdout, stderr string, err error) {
	var stdOutBuf, stdErrBuf bytes.Buffer
	process.Stdout = &stdOutBuf
	process.Stderr = &stdErrBuf
	if err = process.Run(); err != nil {
		return "", "", err
	}
	return stdOutBuf.String(), stdErrBuf.String(), nil
}
