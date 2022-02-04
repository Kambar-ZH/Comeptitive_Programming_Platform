package tools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"site/internal/datastruct"
	"site/internal/logger"
	"site/test/inmemory"
	"time"
)

// Copy content from src to dst
func CopyFile(src, dst string) error {
	file, err := os.Create(dst)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	file.Write(data)
	return nil
}

func BuildExe(execFile, progFile string) error {
	execArg, progArg := fmt.Sprintf("exec='%s'", execFile), fmt.Sprintf("prog='%s'", progFile)
	cmd := exec.Command("make", "run", progArg, execArg)
	cmd.Dir = inmemory.MakeMe()
	var errOut bytes.Buffer
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		logger.Logger.Error(errOut.String())
		return err
	}
	return nil
}

// Run executable file
func ExecuteFile(filePath string, testCase datastruct.TestCase) (string, error) {
	cmd := exec.Command(filePath)

	stdin, err := cmd.StdinPipe()
	defer func() {
		if err := stdin.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		logger.Logger.Error("error on stdinpipe")
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	defer func() {
		if err := stdout.Close(); err != nil {
			return
		}
	}()

	if err != nil {
		return "", err
	}

	if err := cmd.Start(); err != nil {
		return "", err
	}

	input, err := os.ReadFile(testCase.TestFile)
	if err != nil {
		return "", err
	}

	_, err = stdin.Write(input)
	if err != nil {
		return "", err
	}

	result, err := ioutil.ReadAll(stdout)
	if err != nil {
		return "", err
	}
	return string(result), nil
}

// Run executable file, if program doesn't halt, finish after 3 seconds
func MustExecuteFile(filePath string, testCase datastruct.TestCase) (string, error) {
	ch, errorCh := make(chan string), make(chan error)

	go func() {
		verdict, err := ExecuteFile(filePath, testCase)
		if err == nil {
			ch <- verdict
			return
		}
		errorCh <- err
	}()

	select {
	case verdict := <-ch:
		return verdict, nil
	case err := <-errorCh:
		return "", err
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("execution timeout")
	}
}