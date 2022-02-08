package tools

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"site/internal/datastruct"
	"site/internal/logger"
	"site/test/inmemory"
	"time"
)

func ExecuteFile(filePath string, testCase datastruct.TestCase) (string, error) {
	env := fmt.Sprintf("filePath='%s'", filePath)
	cmd := exec.Command("make", "run", env)

	cmd.Dir = inmemory.MakeMe()

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
	type Response struct {
		verdict string
		err     error
	}
	ch := make(chan Response)

	go func() {
		verdict, err := ExecuteFile(filePath, testCase)
		ch <- Response{
			verdict: verdict,
			err:     err,
		}
	}()

	select {
	case response := <-ch:
		return response.verdict, response.err
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("execution timeout")
	}
}
