package main

import (
	"bufio"
	"bytes"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := make(Environment, len(files))

	for _, file := range files {
		val, err := getValueFromFile(dir, file)
		if err != nil {
			return nil, err
		}
		env[file.Name()] = prepareEnvValue(val, file.Size())
	}

	return env, nil
}

func getValueFromFile(currDir string, fileInfo fs.FileInfo) (string, error) {
	pathToFile := currDir + string(filepath.Separator) + fileInfo.Name()
	file, err := os.Open(pathToFile)
	if err != nil {
		return "", err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	val := scanner.Bytes()
	val = bytes.ReplaceAll(val, []byte("\x00"), []byte("\n"))

	strVal := string(val)
	strVal = strings.TrimRight(strVal, " ")

	return strVal, nil
}

func prepareEnvValue(fileValue string, fSize int64) EnvValue {
	if fSize == 0 {
		return EnvValue{"", true}
	}

	return EnvValue{fileValue, false}
}
