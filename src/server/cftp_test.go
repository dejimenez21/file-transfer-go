package main

import (
	"reflect"
	"testing"
)

func TestDeserializeValidCommandWithFileContent(t *testing.T) {
	//given
	inputString := "SEND\n{ \"hasFileContent\": true }\nchn1,chn2\n{ \"name\": \"testDoc\", \"ext\": \"docx\", \"size\": 12017, \"content\": [84, 101, 115, 116, 10, 32, 99, 111, 110, 116, 101, 110, 116, 33, 33, 33] }"
	expectedCommand := command{
		Method:   "SEND",
		Meta:     metaData{},
		Channels: []string{"chn1", "chn2"},
		FileInfo: file{Name: "testDoc", Ext: "docx", Size: 12017},
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
		Meta:     metaData{SenderAddress: "localhost:4567"},
		Channels: []string{"chn1"},
		FileInfo: file{Name: "testDoc", Ext: "docx", Size: 12017},
	}
	expectedBytes := []byte("DELIVER\n{\"SenderAddress\":\"localhost:4567\"}\nchn1\n{\"Name\":\"testDoc\",\"Ext\":\"docx\",\"Size\":12017}\x04")

	//when
	actualBytes, err := serializeCommand(inputCommand)

	//then
	if !reflect.DeepEqual(actualBytes, expectedBytes) || err != nil {
		t.Fatalf("Expected: %v\n got: %v\n error: %v\n", string(expectedBytes), string(actualBytes), err)
	}
}
