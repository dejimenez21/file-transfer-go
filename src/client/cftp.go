package main

import (
	"encoding/json"
	"fmt"
	"strconv"
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

func deserializeRequest(requestString string) (req request, err error) {
	args := strings.SplitN(requestString, "\n", 4)
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

func deserializeDelivery(deliveryString string) (del delivery, err error) {
	del = delivery{}
	deliveryString = strings.TrimSuffix(deliveryString, "\x04")
	args := strings.SplitN(deliveryString, "\n", 4)

	deliveryId, err := strconv.Atoi(args[1])
	if err != nil {
		err = fmt.Errorf("error deserializing chunk: %v", err)
		return
	}
	seq, err := strconv.Atoi(args[2])
	if err != nil {
		err = fmt.Errorf("error deserializing chunk: %v", err)
		return
	}
	size, err := strconv.Atoi(args[3])
	if err != nil {
		err = fmt.Errorf("error deserializing chunk: %v", err)
		return
	}
	del.DeliveryId = int64(deliveryId)
	del.Seq = int64(seq)
	del.Size = size

	return
}
