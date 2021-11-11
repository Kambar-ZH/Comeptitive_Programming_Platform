package services

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUploadFileServiceImpl_PrepareExe(t *testing.T) {
	tests := []struct {
		solutionFile string
		tempFile     string
		foundErr     bool
	}{
		{
			solutionFile: "C:/Users/User/Documents/Visual Studio/GoLang/test/test/problems/0001/solution.go",
			// hardcoded
			tempFile:     "C:/Users/User/Documents/Visual Studio/GoLang/test/temp_solutions/upload-127505235",
			foundErr:     false,
		},
	}
	for _, test := range tests {
		err := PrepareExe(test.solutionFile, test.tempFile)
		var foundErr bool
		if (err != nil) {
			foundErr = true
		}
		assert.Equal(t, test.foundErr, foundErr, fmt.Sprintf("Error on preparing exe. Expected: %v. Found: %v", test.foundErr, foundErr))
	}
}