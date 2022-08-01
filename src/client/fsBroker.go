package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"
)

type fsBroker struct {
	path string
}

func (b *fsBroker) saveFile(f file) (filePath string, err error) {
	if _, err := os.Stat(b.path + f.Name); errors.Is(err, os.ErrNotExist) {
		filePath = b.path + f.Name
	} else {
		filePath = fmt.Sprintf("%s%s_%d", b.path, f.Name, startTime-time.Now().UnixNano())
	}

	newFile, err := os.Create(filePath)
	if err != nil {
		return filePath, fmt.Errorf("an error occurred while creating the file: %v", err)
	}

	newFile.Write(f.Content)
	return
}

func (b *fsBroker) loadFile(path string) (f file, err error) {
	fmt.Println("Getting file", path, "...")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("an error occurred while reading the file: %v", err)
		return
	}
	f.Name = filepath.Base(path)
	f.Ext = filepath.Ext(path)
	f.Content = data
	return
}

func (b *fsBroker) getFileSize(filePath string) (size int64, err error) {
	fileStat, err := os.Stat(filePath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return 0, fmt.Errorf("file %s not found", filePath)
		}
		return 0, fmt.Errorf("error occurred while getting file size: %v", err)
	}
	return fileStat.Size(), nil
}