package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testFilePathFrom = "testdata/input.txt"
	testFilePathTo   = "testdata/result.txt"
)

func TestCopyNegative(t *testing.T) {
	t.Run("file not found", func(t *testing.T) {
		from := "/there/is/no/file"
		to := "/"

		result := Copy(from, to, 0, 0)
		require.ErrorIs(t, result, os.ErrNotExist)
	})

	t.Run("empty paths", func(t *testing.T) {
		from := "/there/is/no/file"
		emptyFrom := ""
		emptyTo := ""

		emptyFromResult := Copy(emptyFrom, testFilePathTo, 0, 0)
		emptyToResult := Copy(from, emptyTo, 0, 0)

		require.ErrorIs(t, emptyFromResult, ErrEmptyPath)
		require.ErrorIs(t, emptyToResult, ErrEmptyPath)
	})

	t.Run("offset more than file size", func(t *testing.T) {
		to := "/"
		var offset int64 = 1024 * 1024 * 1024

		result := Copy(testFilePathFrom, to, offset, 0)
		require.ErrorIs(t, result, ErrOffsetExceedsFileSize)
	})

	t.Run("make new file in restricted directory", func(t *testing.T) {
		to := "/dev/go_test.txt"

		result := Copy(testFilePathFrom, to, 0, 0)
		require.ErrorIs(t, result, os.ErrPermission)
	})
}

func TestCopyPositive(t *testing.T) {
	t.Run("full copy is correct", func(t *testing.T) {
		Copy(testFilePathFrom, testFilePathTo, 0, 0)

		fromStat, _ := os.Stat(testFilePathFrom)
		toStat, _ := os.Stat(testFilePathTo)

		fromSize := fromStat.Size()
		toSize := toStat.Size()
		os.Remove(testFilePathTo)

		require.Equal(t, fromSize, toSize)
	})

	t.Run("partial copy is correct", func(t *testing.T) {
		var limit int64 = 10

		Copy(testFilePathFrom, testFilePathTo, 0, limit)

		toStat, _ := os.Stat(testFilePathTo)
		toSize := toStat.Size()
		os.Remove(testFilePathTo)

		require.Equal(t, toSize, limit)
	})
}
