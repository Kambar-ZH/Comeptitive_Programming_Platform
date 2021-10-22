package main

import (
	"fmt"
	"site/internal/models"
	"site/test"
)

func main() {
	expected, err := test.Run(`C:/Users/User/Documents/Visual Studio/GoLang/test/test/problems/0002/solution.exe`, models.TestCase{
		InputFile: `C:/Users/User/Documents/Visual Studio/GoLang/test/test/problems/0002/tests/2.txt`,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println(expected)
}
