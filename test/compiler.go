package test

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"site/internal/models"
	"time"
)

func WriteFile(inFile, fromFile string) error {
	fo, err := os.Create(inFile)
	if err != nil {
		fmt.Println("Error on opening file")
		fmt.Println(err)
		return err
	}

	defer func() {
		if err := fo.Close(); err != nil {
			fmt.Println("Error on closing file")
			return
		}
	}()

	data, err := os.ReadFile(fromFile)
	if err != nil {
		fmt.Println("File write error")
		return err
	}
	fo.Write(data)
	return nil
}

func MakeExe() error {
	cmd := exec.Command("make")
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		fmt.Println("Could not make file")
		fmt.Println(err)
		return err
	}
	fmt.Println(out.String())
	return nil
}

func RunComplete(file string, tCase models.TestCase) (string, error) {
	cmd := exec.Command(file)

	stdin, err := cmd.StdinPipe()
	defer func() {
		if err := stdin.Close(); err != nil {
			fmt.Println("Error on closing stdin file")
			return
		}
	}()
	if err != nil {
		fmt.Println("StdinPipe error")
		return "", err
	}

	stdout, err := cmd.StdoutPipe()
	defer func() {
		if err := stdout.Close(); err != nil {
			fmt.Println("Error on closing stdout file")
			return
		}
	}()

	if err != nil {
		fmt.Println("StdoutPipe error")
		return "", err
	}

	if err := cmd.Start(); err != nil {
		fmt.Println("Cmd start error")
		return "", err
	}

	input, err := os.ReadFile(tCase.InputFile)
	if err != nil {
		fmt.Println("File write error")
		return "", err
	}
	
	_, err = stdin.Write(input)
	if err != nil {
		fmt.Println("Stdin write error")
		return "", err
	}

	result, err := ioutil.ReadAll(stdout)
	if err != nil {
		fmt.Println("Stdin write error")
		return "", err
	}
	return string(result), nil
}

func Run(file string, tCase models.TestCase) (string, error) {
	ch := make(chan string)
	errCh := make(chan error)
	go func() {
		if verdict, err := RunComplete(file, tCase); err == nil {
			ch <- verdict
		} else {
			errCh <-err
		}
	}()
	select {
	case verdict := (<-ch):
		fmt.Println("MEEEEEE")
		return verdict, nil
	case err := (<-errCh):
		fmt.Println("YOUUUUU")
		return "", err
	case <-time.After(3 * time.Second):
		return "", fmt.Errorf("Timeout") 
	}
}