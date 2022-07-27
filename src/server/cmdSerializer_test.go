package main

import (
	"reflect"
	"testing"
)

func TestDeserializeValidCommand(t *testing.T) {
	//given
	inputString := "SEND\n{ \"hasFile\": true }\nchn1,chn2\n{ \"name\": \"testDoc\", \"ext\": \"docx\", \"size\": 12017 }"
	expectedCommand := command{
		Method:   "SEND",
		Meta:     metaData{HasFile: true},
		Channels: []string{"chn1", "chn2"},
		FileInfo: fileMeta{Name: "testDoc", Ext: "docx", Size: 12017},
	}

	//when
	actualCommand, err := deserializeCommand(inputString)

	//then
	if !reflect.DeepEqual(actualCommand, expectedCommand) || err != nil {
		t.Fatalf("Expected: %v\n got: %v\n error: %v\n", expectedCommand, actualCommand, err)
	}
}
