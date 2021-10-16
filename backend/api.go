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
	tmpl, err := template.ParseFiles("index.html")
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
	// FormFile returns the first file for the given key `myFile`
	// it also returns the FileHeader so we can get the Filename,
	// the Header and the size of the file
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
		return
	}

	// read all of the contents of our uploaded file into a
	// byte array
	fileBytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("Error occured on reading file")
		return
	}
	// write this byte array to our temporary file
	tempFile.Write(fileBytes)
	tempFile.Close()

	problemName := r.FormValue("problemName")
	fileName := strings.ReplaceAll(tempFile.Name(), "\\\\", "/")

	passed, failedTest := test.TestSolution(fileName, problemName)

	// data := ViewData{
	// 	FileName:    handler.Filename,
	// 	ProblemName: problemName,
	// 	Passed:      passed,
	// 	FailedTest:  failedTest,
	// }

	fmt.Fprintf(w, "Successfully uploaded file %s for the problem %s\n", fileName, problemName)

	if passed {
		fmt.Fprintf(w, "All tests passed successfully\n")
	} else {
		fmt.Fprintf(w, "Wrong answer on test %d\n", failedTest)
	}

	// tmpl, err := template.ParseFiles("upload.html")
	// if err != nil {
	// 	fmt.Fprintf(w, "Error occured on loading upload page")
	// 	return
	// }
	// tmpl.Execute(w, data)
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