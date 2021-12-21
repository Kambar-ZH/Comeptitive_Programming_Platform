package tools

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"site/internal/datastruct"
	"time"
)


// Copy to file with path f1Path the content of file located in f2Path
func CopyFile(f1Path, f2Path string) error {
	file, err := os.Create(f1Path)
	if err != nil {
		return err
	}

	defer func() {
		if err := file.Close(); err != nil {
			return
		}
	}()

	data, err := os.ReadFile(f2Path)
	if err != nil {
		return err
	}
	file.Write(data)
	return nil
}

// Run makefile with given command, build executable
func BuildExe(command string) error {
	cmd := exec.Command("make", command)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		log.Println("couldn't build exe")
		return err
	}
	return nil
}

// Run executable file
func ExecuteFile(fPath string, tCase datastruct.TestCase) (string, error) {
	cmd := exec.Command(fPath)

	stdin, err := cmd.StdinPipe()
	defer func() {
		if err := stdin.Close(); err != nil {
			return
		}
	}()
	if err != nil {
		log.Println("error on stdinpipe")
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

	input, err := os.ReadFile(tCase.TestFile)
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
func MustExecuteFile(fPath string, tCase datastruct.TestCase) (string, error) {
	ch, errorCh := make(chan string), make(chan error)

	go func() {
		verdict, err := ExecuteFile(fPath, tCase)
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
		return "", fmt.Errorf("solution timeout")
	}
}
