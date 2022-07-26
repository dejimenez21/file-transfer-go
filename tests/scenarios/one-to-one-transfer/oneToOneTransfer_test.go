package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
)

const SERVER_EXEC_PATH = "../../../src/server/"
const CLIENT_EXEC_PATH = "../../../src/server/"

const CLIENTB_DEFAULT_FOLDER = "./receiverFolder/"

func TestOneToOneTransfer(t *testing.T) {
	//given
	filePath, fileName, fileContent := setupTest()
	expectedFileContent := fileContent
	server := exec.Command("go", "run", SERVER_EXEC_PATH)
	clientA := exec.Command("go", "run", CLIENT_EXEC_PATH, "send", "chn1", filePath+fileName)
	clientB := exec.Command("go", "run", CLIENT_EXEC_PATH, "receive", "chn1", CLIENTB_DEFAULT_FOLDER)

	//when
	server.Run()
	clientB.Run()
	clientA.Run()

	//then
	actualFileContent, err := ioutil.ReadFile(CLIENTB_DEFAULT_FOLDER + fileName)

	if !reflect.DeepEqual(actualFileContent, expectedFileContent) || err != nil {
		t.Fatalf("Error: %v", err)
	}
}

func setupTest() (filePath string, fileName string, fileContent []byte) {
	fileName = "autoTest.txt"
	filePath = "./senderFiles/"
	fileContent = []byte("This is a txt file for testing the one to one transfer!!!\n")
	fileToRemove := fmt.Sprint(CLIENTB_DEFAULT_FOLDER, fileName)
	os.Remove(fileToRemove)
	os.Remove(filePath + fileName)

	file, _ := os.Create(filePath + fileName)
	defer file.Close()

	file.Write(fileContent)
	return
}
