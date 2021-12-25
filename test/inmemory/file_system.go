package inmemory

func RegisterHtml() string {
	relativePath := "web/template/register.html"
	return AbsolutePath(relativePath)
}
func LoginHtml() string {
	relativePath := "web/template/login.html"
	return AbsolutePath(relativePath)
}
func IndexHtml() string {
	relativePath := "web/template/index.html"
	return AbsolutePath(relativePath)
}
func UploadHtml() string {
	relativePath := "web/template/upload.html"
	return AbsolutePath(relativePath)
}
func RatingsHtml() string {
	relativePath := "web/template/ratings.html"
	return AbsolutePath(relativePath)
}
func TempSolutions() string {
	relativePath := "temp_solutions"
	return AbsolutePath(relativePath)
}
func MakeMe() string {
	relativePath := "makeme"
	return AbsolutePath(relativePath)
}
func ParticipantSolution() string {
	relativePath := "cmd/myapp/participant_solution/participant_solution.go"
	return AbsolutePath(relativePath)
}
func ParticipantSolutionExe() string {
	relativePath := "cmd/myapp/participant_solution/participant_solution.exe"
	return AbsolutePath(relativePath)
}
func MainSolution() string {
	relativePath := "cmd/myapp/main_solution/main_solution.go"
	return AbsolutePath(relativePath)
}
func MainSolutionExe() string {
	relativePath := "cmd/myapp/main_solution/main_solution.exe"
	return AbsolutePath(relativePath)
}
func CleanFile() string {
	relativePath := "cmd/myapp/main_solution/clean.txt"
	return AbsolutePath(relativePath)
}

func AbsolutePath(relativePath string) string {
	return "C:/Users/User/Documents/Visual Studio/GoLang/test/" + relativePath
}