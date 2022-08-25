package main

import (
	"log"
	"os"
)

type metaData struct {
	SenderAddress string
	RequestId     int
	Message       string
}

type request struct {
	Method      string
	Meta        metaData
	Channels    []string
	FileInfo    file
	FileContent []byte
}

func (req *request) getMessageType() string {
	return "request"
}

func (req *request) process() (err error) {
	switch req.Method {
	case REQ_DELIVER:
		err = req.processDeliver()
	}
	return
}

func (req *request) processDeliver() (err error) {
	fileBroker := fact.getFileBroker()

	newFile, contentChan, err := fileBroker.saveFile(req.FileInfo, req.Meta.RequestId)
	if err != nil {
		if newFile != nil {
			newFile.Close()
		}
		return err
	}

	go func(contentChan <-chan *delivery, f *os.File) {
		fileBroker := fact.getFileBroker()
		defer f.Close()
		var bytesReceived int64 = 0
		expectedSeq := 0
		for {
			expectedSeq++
			chunk := <-contentChan
			if chunk.Seq != expectedSeq {
				log.Printf("chunk of %s file doesn't have the expected sequence, this indicates an inconsistency in the transmission.", f.Name())
				fileBroker.removeContentChannel(req.Meta.RequestId)
				//TODO: Delete file
			}

			n, err := f.Write(chunk.Content)
			if err != nil {
				log.Printf("error writing content to %s file: %v", f.Name(), err)
				fileBroker.removeContentChannel(req.Meta.RequestId)
				//TODO: Delete file
				return
			}

			bytesReceived += int64(n)
			if bytesReceived >= req.FileInfo.Size {
				log.Printf("successfully received %s file from %s through channel %s", req.FileInfo.Name, req.Meta.SenderAddress, req.Channels[0])
				fileBroker.removeContentChannel(req.Meta.RequestId)
				return
			}
		}
	}(contentChan, newFile)

	log.Printf("started receiving %s file from %s through channel %s", req.FileInfo.Name, req.Meta.SenderAddress, req.Channels[0])

	return
}
