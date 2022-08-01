package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func serializeRequest(req request) (cftpBytes []byte, err error) {
	method := req.Method
	metaBytes, err := json.Marshal(req.Meta)
	if err != nil {
		err = fmt.Errorf("failed to serialize metadata: %v", err)
		return
	}
	meta := string(metaBytes)
	channels := strings.Join(req.Channels, ",")
	fileInfoBytes, err := json.Marshal(req.FileInfo)
	if err != nil {
		err = fmt.Errorf("failed to serialize file information: %v", err)
		return
	}
	fileInfo := string(fileInfoBytes)

	cftpString := strings.Join([]string{method, meta, channels, fileInfo}, "\n")
	cftpBytes = []byte(cftpString)

	return
}

func deserializeRequest(commandString string) (req request, err error) {
	args := strings.SplitN(commandString, "\n", 4)
	method := args[0]
	var meta metaData
	err = json.Unmarshal([]byte(args[1]), &meta)
	if err != nil {
		err = fmt.Errorf("meta section has invalid format: %v", err)
		return
	}
	channels := strings.Split(args[2], ",")
	var finfo file
	err = json.Unmarshal([]byte(args[3]), &finfo)
	if err != nil {
		err = fmt.Errorf("file metadata section has invalid format: %v", err)
		return
	}
	req = request{
		Method:   method,
		Meta:     meta,
		Channels: channels,
		FileInfo: finfo,
	}
	return req, err
}
