package inmemory

func RegisterHtml() string {
	relativePath := "web/template/register.html"
	return absolutePath(relativePath)
}
func LoginHtml() string {
	relativePath := "web/template/login.html"
	return absolutePath(relativePath)
}
func IndexHtml() string {
	relativePath := "web/template/index.html"
	return absolutePath(relativePath)
}
func UploadHtml() string {
	relativePath := "web/template/upload.html"
	return absolutePath(relativePath)
}
func RatingsHtml() string {
	relativePath := "web/template/ratings.html"
	return absolutePath(relativePath)
}
func TempSolutions() string {
	relativePath := "temp_solutions"
	return absolutePath(relativePath)
}
func ParticipantSolution() string {
	relativePath := "cmd/myapp/participant_solution/participant_solution.go"
	return absolutePath(relativePath)
}
func ParticipantSolutionExe() string {
	relativePath := "cmd/myapp/participant_solution/participant_solution.exe"
	return absolutePath(relativePath)
}
func MainSolution() string {
	relativePath := "cmd/myapp/main_solution/main_solution.go"
	return absolutePath(relativePath)
}
func MainSolutionExe() string {
	relativePath := "cmd/myapp/main_solution/main_solution.exe"
	return absolutePath(relativePath)
}
func CleanFile() string {
	relativePath := "cmd/myapp/main_solution/clean.txt"
	return absolutePath(relativePath)
}

func absolutePath(relativePath string) string {
	return "C:/Users/User/Documents/Visual Studio/GoLang/test/" + relativePath
}