package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

type fsBroker struct {
	path             string
	contentChans     map[int]chan *delivery
	contentChansLock sync.RWMutex
}

func (b *fsBroker) saveFile(f file, deliveryId int) (newFile *os.File, contentChan <-chan *delivery, err error) {
	b.createFolder()
	filePath := b.getSavingFileFullName(f)

	newFile, err = os.Create(filePath)
	if err != nil {
		err = fmt.Errorf("an error occurred while creating the file: %v", err)
		return
	}
	contentChan = b.newDeliveryChannel(deliveryId)

	return
}

func (b *fsBroker) newDeliveryChannel(deliveryId int) <-chan *delivery {
	b.contentChansLock.Lock()
	defer b.contentChansLock.Unlock()
	if b.contentChans == nil {
		b.contentChans = make(map[int]chan *delivery)
	}
	newChan := make(chan *delivery)
	b.contentChans[deliveryId] = newChan
	return newChan
}

func (b *fsBroker) removeContentChannel(deliveryId int) {
	b.contentChansLock.Lock()
	defer b.contentChansLock.Unlock()
	c, found := b.contentChans[deliveryId]
	if !found {
		return
	}

	select {
	case <-c:
		delete(b.contentChans, deliveryId)
	default:
		delete(b.contentChans, deliveryId)
	}
}

func (b *fsBroker) getSavingFileFullName(f file) (filePath string) {
	if _, err := os.Stat(b.path + f.Name); errors.Is(err, os.ErrNotExist) {
		filePath = b.path + f.Name
	} else {
		name := strings.TrimSuffix(f.Name, f.Ext)
		filePath = fmt.Sprintf("%s%s_%d%s", b.path, name, startTime-time.Now().UnixNano(), f.Ext)
	}
	return
}

func (b *fsBroker) loadFile(path string) (f file, content []byte, err error) {
	fmt.Println("Getting file", path, "...")
	data, err := ioutil.ReadFile(path)
	if err != nil {
		err = fmt.Errorf("an error occurred while reading the file: %v", err)
		return
	}
	f.Name = filepath.Base(path)
	f.Ext = filepath.Ext(path)
	content = data
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

func (b *fsBroker) createFolder() {
	os.MkdirAll(b.path, 0700)
}

func (b *fsBroker) saveChunk(chunk *delivery) (err error) {
	b.contentChansLock.RLock()
	defer b.contentChansLock.RUnlock()

	contentChan, found := b.contentChans[chunk.DeliveryId]
	if !found {
		err = fmt.Errorf("unexpected deliveryId")
		return
	}

	contentChan <- chunk

	return
}
