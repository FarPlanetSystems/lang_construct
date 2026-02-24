package main

import (
	"fmt"
	"os"
	"strings"
)

func readCodeFile(file_path string) (string, error) {
	_, err := os.Open(file_path)
	if err != nil {
		return "", err
	}
	bytes, _ := os.ReadFile(file_path)
	code := ""
	for i := 0; i < len(bytes); i++ {
		if string(bytes[i]) != "@" {
			code += string(bytes[i])
		} else {
			return code, nil
		}
	}
	return code, nil
}

func ReadCode() string {
	// we get the file name of a .txt from the terminal
	file_path, err := getFilePath()
	// check that everything was alright by getting file name
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// get the code inside the file
	code, err := readCodeFile(file_path)
	// check that everything was alright opening and reading the file
	if err != nil {
		fmt.Println(err)
		return ""
	}
	// check whether the file is empty
	if len(code) < 1 {
		fmt.Println("cannot run an empty file " + file_path)
		return ""
	}
	return code
}

func getFilePath() (string, error) {
	file_name := os.Args[1]
	// parsing
	file_name = strings.ReplaceAll(file_name, " ", "")
	file_name = strings.ReplaceAll(file_name, "\n", "")
	file_name = strings.ReplaceAll(file_name, "\r", "")
	// creating the path
	curdir, err := os.Getwd()
	file_path := curdir + "\\" + file_name
	return file_path, err
}
