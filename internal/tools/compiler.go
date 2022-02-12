package tools

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"regexp"
	"site/internal/datastruct"
	"site/internal/logger"
	"site/test/inmemory"
	"strings"
	"time"
)

var (
	ExtensionRegex = regexp.MustCompile(`.*\.`)
)

func ExecuteFile(filePath string, testCase datastruct.TestCase) (string, error) {
	extension := ExtensionRegex.ReplaceAllString(filePath, "")
	
	var cmd *exec.Cmd

	switch extension {

	case "go":
		env := fmt.Sprintf("filePath='%s'", filePath)
		cmd = exec.Command("make", "run_go_file", env)

	case "c++":
		env := fmt.Sprintf("filePath='%s'", filePath)
		executablePath := strings.TrimSuffix(filePath, extension) + "out"
		env2 := fmt.Sprintf("executablePath=%s", executablePath)
		cmd = exec.Command("make", "run_cpp_file", env, env2)
		defer func() {
			cmd = exec.Command("make", "rm_cpp_executable", env2)
			cmd.Dir = inmemory.MakeMe()

			var errorOut bytes.Buffer
			cmd.Stderr = &errorOut
			
			if err := cmd.Run(); err != nil {
				logger.Logger.Error(errorOut.String())
			}
		}()

	case "py":
		env := fmt.Sprintf("filePath='%s'", filePath)
		cmd = exec.Command("make", "run_py_file", env)
	}

	cmd.Dir = inmemory.MakeMe()

	var errorOut bytes.Buffer
	cmd.Stderr = &errorOut

	stdin, err := cmd.StdinPipe()
	defer func() {
		if err := stdin.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		logger.Logger.Error("error on stdinpipe")
		return "", errors.New(errorOut.String())
	}

	stdout, err := cmd.StdoutPipe()
	defer func() {
		if err := stdout.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		return "", errors.New(errorOut.String())
	}

	if err := cmd.Start(); err != nil {
		return "", errors.New(errorOut.String())
	}

	input, err := os.ReadFile(testCase.TestFile)
	if err != nil {
		return "", err
	}

	_, err = stdin.Write(input)
	if err != nil {
		return "", errors.New(errorOut.String())
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
