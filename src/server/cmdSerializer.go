package main

import (
	"encoding/json"
	"fmt"
	"strings"
)

func deserializeCommand(commandString string) (cmd command, err error) {
	args := strings.Split(commandString, "\n")
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
