package test

import (
	"fmt"
)

type TestCase struct {
	inputFile string
}

type Validator struct {
	ProblemName  string
	MainSolution string
	TestCases    []TestCase
}

var (
	tests = []Validator {
		{
			"A",
			"../test/problems/A/solution.exe",
			[]TestCase{
				{
					"../test/problems/A/tests/1.txt",
				},
				{
					"../test/problems/A/tests/2.txt",
				},
			},
		},
		{
			"B",
			"../test/problems/B/solution.exe",
			[]TestCase{
				{
					"../test/problems/B/tests/1.txt",
				},
				{
					"../test/problems/B/tests/2.txt",
				},
			},
		},
	}
)

func TestSolution(participantSolution string, problemName string) (bool, int) {
	for _, test := range tests {
		if test.ProblemName == problemName {
			for tId, tCase := range test.TestCases {
				expected := Run(test.MainSolution, tCase)
				actual := Run(participantSolution, tCase)
				if expected != actual {
					fmt.Printf("Fail.\n\tExpected: %s\n\tActual: %s\n", expected, actual)
					return false, tId + 1
				}
				fmt.Printf("OK.\n\tExpected: %s\n\tActual: %s\n", expected, actual)
			}
		}
	}
	return true, -1
}