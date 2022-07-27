package main

import (
	"reflect"
	"testing"
)

func TestDeserializeValidCommandWithFileContent(t *testing.T) {
	//given
	inputString := "SEND\n{ \"hasFile\": true }\nchn1,chn2\n{ \"name\": \"testDoc\", \"ext\": \"docx\", \"size\": 12017, \"content\": [84, 101, 115, 116, 10, 32, 99, 111, 110, 116, 101, 110, 116, 33, 33, 33] }"
	expectedCommand := command{
		Method:   "SEND",
		Meta:     metaData{HasFile: true},
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
