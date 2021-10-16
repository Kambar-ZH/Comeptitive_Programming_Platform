package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"site/test"
	"strings"
	"text/template"

	"github.com/gorilla/mux"
)

type ViewData struct {
	FileName    string
	ProblemName string
	Passed      bool
	FailedTest  int
}

func homePage(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("../frontend/index.html")
	if err != nil {
		fmt.Fprintf(w, "Error occured on loading home page")
		fmt.Println(err)
		return
	}
	data := ViewData{}
	tmpl.Execute(w, data)
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Println("File Upload Endpoint Hit")

	r.ParseMultipartForm(10 << 20)

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		return
	}

	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// Create a temporary file within our temp-images directory that follows
	// a particular naming pattern
	tempFile, err := ioutil.TempFile("../temp-solutions", "upload-*.exe")
	if err != nil {
		fmt.Println("Error occured on creating temp file")
		fmt.Println(err)
		return
	}

	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error occured on reading file")
		fmt.Println(err)
		return
	}
	
	tempFile.Write(fileBytes)
	tempFile.Close()

	problemName := r.FormValue("problemName")
	fileName := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")

	passed, failedTest := test.TestSolution(fileName, problemName)

	data := ViewData{
		FileName:    handler.Filename,
		ProblemName: problemName,
		Passed:      passed,
		FailedTest:  failedTest,
	}

	tmpl, err := template.ParseFiles("../frontend/upload.html")
	if err != nil {
		fmt.Fprintf(w, "Error occured on loading upload page")
		fmt.Println(err)
		return
	}
	tmpl.Execute(w, data)
}

func setupRoutes() {
	r := mux.NewRouter()
	r.HandleFunc("/upload", uploadFile)
	r.HandleFunc("/", homePage)
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.ListenAndServe(":8080", r)
}

func main() {
	fmt.Println("Server started work")
	setupRoutes()
	fmt.Println("Server ends work")
}