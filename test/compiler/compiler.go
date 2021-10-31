package compiler

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"site/internal/models"
	"time"
)

func CopyFile(inFile, fromFile string) error {
	fo, err := os.Create(inFile)
	if err != nil {
		log.Println("error on opening file")
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			log.Println("error on closing file")
			return
		}
	}()

	data, err := os.ReadFile(fromFile)
	if err != nil {
		log.Println("error on reading file")
		return err
	}
	fo.Write(data)
	return nil
}

func BuildExe() error {
	cmd := exec.Command("make")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		log.Println("error: could not make file")
		return err
	}
	return nil
}

func ExecuteFile(file string, tCase models.TestCase) (string, error) {
	cmd := exec.Command(file)

	stdin, err := cmd.StdinPipe()
	defer func() {
		if err := stdin.Close(); err != nil {
			log.Println("error on closing stdin file")
			return
		}
	}()
	if err != nil {
		log.Println("stdinPipe error")
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	defer func() {
		if err := stdout.Close(); err != nil {
			log.Println("error on closing stdout file")
			return
		}
	}()

	if err != nil {
		log.Println("stdoutPipe error")
		return "", err
	}

	if err := cmd.Start(); err != nil {
		log.Println("cmd start error")
		return "", err
	}

	input, err := os.ReadFile(tCase.InputFile)
	if err != nil {
		log.Println("file write error")
		return "", err
	}

	_, err = stdin.Write(input)
	if err != nil {
		log.Println("stdin write error")
		return "", err
	}

	result, err := ioutil.ReadAll(stdout)
	if err != nil {
		log.Println("stdin write error")
		return "", err
	}
	return string(result), nil
}

func MustExecuteFile(file string, tCase models.TestCase) (string, error) {
	ch := make(chan string)
	errorCh := make(chan error)

	go func() {
		if verdict, err := ExecuteFile(file, tCase); err == nil {
			ch <- verdict
		} else {
			errorCh <- err
		}
	}()

	select {
	case verdict := <-ch:
		return verdict, nil
	case err := <-errorCh:
		return "", err
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("server timeout")
	}
}