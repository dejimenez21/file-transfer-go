package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func deserializeCommand(commandString string) (cmd command, err error) {
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
	cmd = command{
		Method:   method,
		Meta:     meta,
		Channels: channels,
		FileInfo: finfo,
	}
	return cmd, err
}

func serializeCommand(cmd command) (cftpBytes []byte, err error) {
	method := cmd.Method
	metaBytes, err := json.Marshal(cmd.Meta)
	if err != nil {
		err = fmt.Errorf("failed to serialize metadata: %v", err)
		return
	}
	meta := string(metaBytes)
	channels := strings.Join(cmd.Channels, ",")
	fileInfoBytes, err := json.Marshal(cmd.FileInfo)
	if err != nil {
		err = fmt.Errorf("failed to serialize file information: %v", err)
		return
	}
	fileInfo := string(fileInfoBytes)

	cftpString := strings.Join([]string{method, meta, channels, fileInfo}, "\n")
	cftpBytes = []byte(cftpString)
	cftpBytes = append(cftpBytes, EOT)
	return
}

func serializeChunkDelivery(del delivery) (bytes []byte) {
	method := "chunk"
	dataString := strings.Join([]string{method, fmt.Sprint(del.ID), fmt.Sprint(del.Seq), fmt.Sprint(del.Size)}, "\n")
	dataString += "\x04"
	bytes = append([]byte(dataString), del.Content...)
	return
}
