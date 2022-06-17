package main

import (
	"bytes"
	"fmt"
	"github.com/Umesh-Kumar-01/autograder/logs"
	"os"
	"os/exec"
	"path/filepath"
	"sync"
)

func addDetailsInTextFile(username string, testPass bool) {
	file, err := os.OpenFile("result.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		logs.ErrorLogger.Println(err)
		return
	}
	defer file.Close()
	var toWrite string
	passOrNot := "PASS"
	if !testPass {
		passOrNot = "FAIL"
	}
	toWrite = fmt.Sprintln(username, passOrNot)
	file.WriteString(toWrite)
}

func fileEndingWithGivenString(path string, endingWith string) bool {
	// check if path contains "_test.go" in the last of file name
	var sz int = len(path)
	var szReq int = len(endingWith) // size of path should be at least > sizeof(_test.go)
	if sz > szReq {
		if path[sz-szReq:sz] == endingWith {
			return true
		}
		return false
	}
	return false
}

func findGoFilesForTest(path string) (isGoFile bool, fileName string, testFileName string) {
	if fileEndingWithGivenString(path, ".go") && !fileEndingWithGivenString(path, "_test.go") {
		isGoFile = true
		fileName = path
		testFileName = path[:len(path)-3] + "_test.go"
		return
	}
	return false, "", ""
}
func testGoFile(path string) {
	x, fileName, testFileName := findGoFilesForTest(path)
	if !x {
		logs.WarningLogger.Println("This file ->", path, "is not a go file or it is a go test file.")
		return
	}
	out := bytes.Buffer{}
	cmd := exec.Command("go", "test", fileName, testFileName)
	cmd.Stdout = &out
	err := cmd.Run()
	if err != nil {
		logs.ErrorLogger.Println("Error for", path, "Error:", err)
		return
	}
	checker := make([]byte, 2)
	n, er := out.Read(checker)
	if er != nil || n != 2 {
		logs.ErrorLogger.Println("Error for", path, "Error:", err)
		return
	}
	if checker[0] == 'o' && checker[1] == 'k' {
		addDetailsInTextFile(path, true)
		return
	}
	addDetailsInTextFile(path, false)
}

func main() {
	var root string
	fmt.Scanf("Write the absolute path of files where all submissions are present:\n&s", &root)
	var files []string

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		files = append(files, path)
		return nil
	})
	if err != nil {
		logs.ErrorLogger.Println(err)
		return
	}
	w := sync.WaitGroup{}
	for _, file := range files {
		func() {
			w.Add(1)
			go testGoFile(file)
			w.Done()
		}()
	}
	w.Wait()
}
