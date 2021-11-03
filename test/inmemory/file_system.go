package inmemory

import (
	"fmt"
	"site/internal/models"
)

var (
	fileSystem *FileSystem = nil
)

type FileSystem struct {
	Tests                  map[int]models.Validator
}

func (fs *FileSystem) GetTestByID(id int) (*models.Validator, error) {
	validator, ok := fs.Tests[id]
	if !ok {
		return nil, fmt.Errorf("problem with this id not found")
	}
	return &validator, nil
}
func (fs *FileSystem) GetRegisterHtml() string {
	relativePath := "web/template/register.html"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetLoginHtml() string {
	relativePath := "web/template/login.html"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetIndexHtml() string {
	relativePath := "web/template/index.html"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetUploadHtml() string {
	relativePath := "web/template/upload.html"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetRatingsHtml() string {
	relativePath := "web/template/ratings.html"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetTempSolutions() string {
	relativePath := "temp_solutions"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetParticipantSolution() string {
	relativePath := "cmd/myapp/participant_solution/participant_solution.go"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetParticipantSolutionExe() string {
	relativePath := "cmd/myapp/participant_solution/participant_solution.exe"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetMainSolution() string {
	relativePath := "cmd/myapp/main_solution/main_solution.go"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetMainSolutionExe() string {
	relativePath := "cmd/myapp/main_solution/main_solution.exe"
	return absolutePath(relativePath)
}
func (fs *FileSystem) GetCleanFile() string {
	relativePath := "cmd/myapp/main_solution/clean.txt"
	return absolutePath(relativePath)
}

func absolutePath(relativePath string) string {
	return "C:/Users/User/Documents/Visual Studio/GoLang/test/" + relativePath
}

func GetInstance() *FileSystem {
	if fileSystem == nil {
		fileSystem = &FileSystem{
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
	return fileSystem
}
