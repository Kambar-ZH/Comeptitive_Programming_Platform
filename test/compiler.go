package test

import (
	"io/ioutil"
	"os"
	"os/exec"
)

func ReadFile(fileName string) []byte {
	file, err := os.ReadFile(fileName)
	if err != nil {
		panic(err)
	}
	return file
}

func Run(file string, tCase TestCase) string {
	cmd := exec.Command(file)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		panic(err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		panic(err)
	}

	if err := cmd.Start(); err != nil {
		panic(err)
	}

	input := ReadFile(tCase.inputFile)
	_, err = stdin.Write(input)
	if err != nil {
		panic(err)
	}
	stdin.Close()

	result, err := ioutil.ReadAll(stdout)
	if err != nil {
		panic(err)
	}
	return string(result)
}