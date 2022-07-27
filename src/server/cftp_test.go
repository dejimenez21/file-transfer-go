package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"reflect"
	"testing"
)

func TestDeserializeValidCommandWithFileContent(t *testing.T) {
	//given
	inputString := "SEND\n{ \"hasFileContent\": true }\nchn1,chn2\n{ \"name\": \"testDoc\", \"ext\": \"docx\", \"size\": 12017, \"content\": [84, 101, 115, 116, 10, 32, 99, 111, 110, 116, 101, 110, 116, 33, 33, 33] }"
	expectedCommand := command{
		Method:   "SEND",
		Meta:     metaData{HasFileContent: true},
		Channels: []string{"chn1", "chn2"},
		FileInfo: file{Name: "testDoc", Ext: "docx", Size: 12017, Content: []byte{84, 101, 115, 116, 10, 32, 99, 111, 110, 116, 101, 110, 116, 33, 33, 33}},
	}

	//when
	actualCommand, err := deserializeCommand(inputString)

	//then
	if !reflect.DeepEqual(actualCommand, expectedCommand) || err != nil {
		t.Fatalf("Expected: %v\n got: %v\n error: %v\n", expectedCommand, actualCommand, err)
	}
}

func TestSerializeValidDeliveryWithFileContent(t *testing.T) {
	//given
	inputCommand := command{
		Method:   "DELIVER",
		Meta:     metaData{HasFileContent: true, SenderAddress: "localhost:4567"},
		Channels: []string{"chn1"},
		FileInfo: file{Name: "testDoc", Ext: "docx", Size: 12017, Content: []byte{84, 101, 115, 116, 10, 32, 99, 111, 110, 116, 101, 110, 116, 33, 33, 33}},
	}
	expectedFileContent := base64.StdEncoding.EncodeToString(inputCommand.FileInfo.Content)
	expectedBytes := []byte(fmt.Sprintf("DELIVER\n{\"HasFileContent\":true,\"SenderAddress\":\"localhost:4567\"}\nchn1\n{\"Name\":\"testDoc\",\"Ext\":\"docx\",\"Size\":12017,\"Content\":\"%s\"}", expectedFileContent))

	//when
	actualBytes, err := serializeDelivery(inputCommand)

	//then
	if !reflect.DeepEqual(actualBytes, expectedBytes) || err != nil {
		t.Fatalf("Expected: %v\n got: %v\n error: %v\n", string(expectedBytes), string(actualBytes), err)
	}
}

func TestJson(t *testing.T) {
	inputCommand := command{
		Method:   "DELIVER",
		Meta:     metaData{HasFileContent: true, SenderAddress: "localhost:4567"},
		Channels: []string{"chn1"},
		FileInfo: file{Name: "testDoc", Ext: "docx", Size: 12017, Content: []byte{84, 101, 115, 116, 10, 32, 99, 111, 110, 116, 101, 110, 116, 33, 33, 33}},
	}

	fileBytes, _ := json.Marshal(inputCommand.FileInfo)
	var actualCommand file
	json.Unmarshal(fileBytes, &actualCommand)

	fmt.Println(actualCommand)
}
