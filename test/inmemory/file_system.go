package inmemory

import (
	"path/filepath"
	"runtime"
	"strings"
)

var (
	absPath string
)

func RegisterHtml() string {
	relativePath := "web/template/register.html"
	return AbsPath(relativePath)
}
func LoginHtml() string {
	relativePath := "web/template/login.html"
	return AbsPath(relativePath)
}
func IndexHtml() string {
	relativePath := "web/template/index.html"
	return AbsPath(relativePath)
}
func UploadHtml() string {
	relativePath := "web/template/upload.html"
	return AbsPath(relativePath)
}
func RatingsHtml() string {
	relativePath := "web/template/ratings.html"
	return AbsPath(relativePath)
}
func TempSolutions() string {
	relativePath := "../../temp_solutions"
	return relativePath
}
func MakeMe() string {
	relativePath := "makeme"
	return AbsPath(relativePath)
}
func AbsPath(relativePath string) string {
	if (absPath == "") {
		_, b, _, _ := runtime.Caller(0)
		basepath := filepath.Dir(b)
		path := strings.Split(basepath, "/")
		path = path[:len(path)-2]
		absPath = strings.Join(path, "/") + "/"
	}
	return absPath + relativePath
}