package main

import (
	"log"
	"os/exec"
)

func main() {
	cmd := exec.Command("make", "run", "prog='1.go'", "exec='1.exe'")
	cmd.Dir = "C:/Users/User/Documents/Visual Studio/GoLang/test/cmd/testapp/makeme"
	err := cmd.Run()
	if err != nil {
		log.Fatal(err)
	}
}