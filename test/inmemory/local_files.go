package inmemory

import (
	"fmt"
	"site/internal/models"
)

var (
	localFiles *LocalFiles = nil
)

type LocalFiles struct {
	IndexHtml              string
	UploadHtml             string
	TempSolutions          string
	ParticipantSolution    string
	ParticipantSolutionExe string
	MainSolution           string
	MainSolutionExe        string
	CleanFile              string
	Tests                  map[int]models.Validator
}

func (lf *LocalFiles) GetTestByID(id int) (*models.Validator, error) {
	validator, ok := localFiles.Tests[id]
	if !ok {
		return nil, fmt.Errorf("the problem with this id not found")
	}
	return &validator, nil
}

func GetInstance() *LocalFiles {
	if localFiles == nil {
		localFiles = &LocalFiles{
			IndexHtml:              "../../web/template/index.html",
			UploadHtml:             "../../web/template/upload.html",
			TempSolutions:          "../../temp_solutions",
			ParticipantSolution:    "../../cmd/myapp/participant_solution/participant_solution.go",
			ParticipantSolutionExe: "../../cmd/myapp/participant_solution/participant_solution.exe",
			MainSolution:           "../../cmd/myapp/main_solution/main_solution.go",
			MainSolutionExe:        "../../cmd/myapp/main_solution/main_solution.exe",
			CleanFile:              "../../cmd/myapp/main_solution/clean.txt",
			Tests: map[int]models.Validator{
				1: {
					ProblemId:    1,
					SolutionFile: "../../test/problems/0001/solution.go",
					TestCases: []models.TestCase{
						{
							InputFile: "../../test/problems/0001/tests/1.txt",
						},
						{
							InputFile: "../../test/problems/0001/tests/2.txt",
						},
					},
				},
				2: {
					ProblemId:    2,
					SolutionFile: "../../test/problems/0002/solution.go",
					TestCases: []models.TestCase{
						{
							InputFile: "../../test/problems/0002/tests/1.txt",
						},
						{
							InputFile: "../../test/problems/0002/tests/2.txt",
						},
					},
				},
			},
		}
	}
	return localFiles
}
