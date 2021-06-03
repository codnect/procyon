package main

import (
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
)

func createFile(filePath, content string) error {
	return ioutil.WriteFile(filePath, []byte(content), os.ModePerm)
}

func checkIfExist(path string) bool {
	_, err := os.Stat(path)
	return !os.IsNotExist(err)
}

func checkIfGoInstalled() error {
	err := exec.Command("go", "help").Run()

	if err != nil {
		return errors.New("please check if golang is installed successfully on your system")
	}

	return nil
}
