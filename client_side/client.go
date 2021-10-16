package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
)

func main() {
    reqBody, err := json.Marshal(map[string]string{
        "code": "HELLO",
    })
	// `package main
	// 	func Solution() string {
	// 		return "Hello"
	// 	}`
    if err != nil {
        print(err)
    }
    resp, err := http.Post("http://localhost:8080/",
        "application/json", bytes.NewBuffer(reqBody))
    if err != nil {
        print(err)
    }
    defer resp.Body.Close()
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        print(err)
    }
    fmt.Println(string(body))
}