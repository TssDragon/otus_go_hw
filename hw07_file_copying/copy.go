package main

import (
	"errors"
	"io"
	"os"

	"github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
	ErrEmptyPath             = errors.New("from/to path must be set")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	err := checkPathsError(fromPath, toPath)
	if err != nil {
		return err
	}

	err = checkFileSizeAndOffset(fromPath, offset)
	if err != nil {
		return err
	}

	currFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	defer closeFileAndProgressBarIfExists(currFile, nil)
	if err != nil {
		return err
	}

	_, err = currFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	bytesForCopy := setRealBytesForCopy(currFile, limit, offset)

	err = realCopy(currFile, toPath, bytesForCopy)
	if err != nil {
		return err
	}

	return nil
}

func checkPathsError(paths ...string) error {
	for _, path := range paths {
		if path == "" {
			return ErrEmptyPath
		}
	}

	return nil
}

func checkFileSizeAndOffset(path string, offset int64) error {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return err
	}

	if fileInfo.Size() == 0 {
		return ErrUnsupportedFile
	}

	if offset > fileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	return nil
}

func setRealBytesForCopy(file *os.File, limit int64, offset int64) int64 {
	fileStat, _ := file.Stat()
	fileSize := fileStat.Size()

	bytesForCopy := fileSize - offset

	if limit < bytesForCopy && limit > 0 {
		bytesForCopy = limit
	}

	return bytesForCopy
}

func realCopy(from *os.File, path string, limit int64) error {
	newFile, err := os.Create(path)
	progressBar := pb.Full.Start64(limit)

	defer closeFileAndProgressBarIfExists(newFile, progressBar)
	if err != nil {
		return err
	}

	proxyReader := progressBar.NewProxyReader(from)

	_, err = io.CopyN(newFile, proxyReader, limit)
	if err != nil {
		return err
	}

	progressBar.Finish()

	return nil
}

func closeFileAndProgressBarIfExists(file *os.File, bar *pb.ProgressBar) {
	if file != nil {
		file.Close()
	}

	if bar != nil {
		bar.Finish()
	}
}
