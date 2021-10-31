package inmemory

import (
	"fmt"
	"site/internal/models"
)

var (
	filesSystem *FileSystem = nil
)

type FileSystem struct {
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

func (lf *FileSystem) GetTestByID(id int) (*models.Validator, error) {
	validator, ok := filesSystem.Tests[id]
	if !ok {
		return nil, fmt.Errorf("problem with this id not found")
	}
	return &validator, nil
}

func absolutePath(relativePath string) string {
	return "C:/Users/User/Documents/Visual Studio/GoLang/test/" + relativePath
}

func GetInstance() *FileSystem {
	if filesSystem == nil {
		filesSystem = &FileSystem{
			IndexHtml:              absolutePath("web/template/index.html"),
			UploadHtml:             absolutePath("web/template/upload.html"),
			TempSolutions:          absolutePath("temp_solutions"),
			ParticipantSolution:    absolutePath("cmd/myapp/participant_solution/participant_solution.go"),
			ParticipantSolutionExe: absolutePath("cmd/myapp/participant_solution/participant_solution.exe"),
			MainSolution:           absolutePath("cmd/myapp/main_solution/main_solution.go"),
			MainSolutionExe:        absolutePath("cmd/myapp/main_solution/main_solution.exe"),
			CleanFile:              absolutePath("cmd/myapp/main_solution/clean.txt"),
			Tests: map[int]models.Validator{
				1: {
					ProblemId:    1,
					SolutionFile: absolutePath("test/problems/0001/solution.go"),
					TestCases: []models.TestCase{
						{
							InputFile: absolutePath("test/problems/0001/tests/1.txt"),
						},
						{
							InputFile: absolutePath("test/problems/0001/tests/2.txt"),
						},
					},
				},
				2: {
					ProblemId:    2,
					SolutionFile: absolutePath("test/problems/0002/solution.go"),
					TestCases: []models.TestCase{
						{
							InputFile: absolutePath("test/problems/0002/tests/1.txt"),
						},
						{
							InputFile: absolutePath("test/problems/0002/tests/2.txt"),
						},
					},
				},
			},
		}
	}
	return filesSystem
}
