package main

import (
	"log"
	"site/internal/datastruct"
	"site/internal/tools"
	"site/test/inmemory"
)

func main() {
	testCase := datastruct.TestCase{
		TestFile: inmemory.AbsPath("test/problems/0001/tests/1.txt"),
	}
	res, err := tools.ExecuteFile(inmemory.AbsPath("test/problems/0001/solution.go"), testCase)
	if err != nil {
		log.Fatal(err)
	}
	print(res)
}