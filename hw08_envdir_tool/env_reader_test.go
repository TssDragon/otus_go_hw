package main

import (
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

const testDir string = "testdata/env"

func TestReadDir(t *testing.T) {
	t.Run("prepare correct env", func(t *testing.T) {
		result, _ := ReadDir(testDir)

		expectedEnv := Environment{
			"BAR":   {"bar", false},
			"EMPTY": {"", false},
			"FOO":   {"   foo\nwith new line", false},
			"HELLO": {"\"hello\"", false},
			"UNSET": {"", true},
		}

		require.Equal(t, expectedEnv, result)
	})

	t.Run("non-existent dir", func(t *testing.T) {
		result, err := ReadDir(testDir + string(filepath.Separator) + "non_existing_dir")
		require.Nil(t, result)
		require.Error(t, err)
	})

	t.Run("empty dir", func(t *testing.T) {
		emptyDirPath, _ := os.MkdirTemp(testDir, "test")
		fmt.Println(emptyDirPath)
		result, _ := ReadDir(emptyDirPath)
		expectedEnv := Environment{}

		require.Equal(t, expectedEnv, result)
	})
}

func TestGetValueFromFile(t *testing.T) {
	t.Run("empty file", func(t *testing.T) {
		pathToFile := testDir + string(filepath.Separator) + "EMPTY"
		fileInfo, _ := os.Stat(pathToFile)

		result, err := getValueFromFile(testDir, fileInfo)
		require.NoError(t, err)
		require.Equal(t, "", result)
	})

	t.Run("single line", func(t *testing.T) {
		pathToFile := testDir + string(filepath.Separator) + "HELLO"
		fileInfo, _ := os.Stat(pathToFile)

		result, err := getValueFromFile(testDir, fileInfo)
		require.NoError(t, err)
		require.Equal(t, "\"hello\"", result)
	})

	t.Run("more than one line", func(t *testing.T) {
		pathToFile := testDir + string(filepath.Separator) + "BAR"
		fileInfo, _ := os.Stat(pathToFile)

		result, err := getValueFromFile(testDir, fileInfo)
		require.NoError(t, err)
		require.Equal(t, "bar", result)
	})

	t.Run("line with terminal nul", func(t *testing.T) {
		pathToFile := testDir + string(filepath.Separator) + "FOO"
		fileInfo, _ := os.Stat(pathToFile)

		result, err := getValueFromFile(testDir, fileInfo)
		expected := "   foo\nwith new line"
		require.NoError(t, err)
		require.Equal(t, expected, result)
	})
}

func TestPrepareEnvValue(t *testing.T) {
	t.Run("prepare env value empty file", func(t *testing.T) {
		result := prepareEnvValue("test", 0)
		require.Equal(t, EnvValue{"", true}, result)
	})

	t.Run("prepare env value non empty file", func(t *testing.T) {
		fVal := "test"

		result := prepareEnvValue(fVal, 10)
		require.Equal(t, EnvValue{fVal, false}, result)
	})
}
