package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"reflect"
	"testing"
	"time"
)

const (
	SERVER_EXEC_PATH       = "C:\\Users\\dejim\\source\\repos\\file-transfer-go\\src\\server\\server.exe"
	CLIENT_EXEC_PATH       = "C:\\Users\\dejim\\source\\repos\\file-transfer-go\\src\\client\\client.exe"
	CLIENTB_DEFAULT_FOLDER = "C:\\Users\\dejim\\source\\repos\\file-transfer-go\\tests\\scenarios\\one-to-one-transfer\\receiverFolder\\"
)

func TestOneToOneTransfer(t *testing.T) {
	//given
	filePath, fileName, fileContent := setupTest()
	expectedFileContent := fileContent
	server := exec.Command(SERVER_EXEC_PATH)
	clientA := exec.Command(CLIENT_EXEC_PATH, "send", "-ch", "chn1", filePath+fileName)
	clientB := exec.Command(CLIENT_EXEC_PATH, "receive", "-ch", "chn1", "-path", CLIENTB_DEFAULT_FOLDER)

	//when
	go server.Run()
	go clientB.Run()
	go clientA.Run()
	time.Sleep(2000000000)

	//then
	actualFileContent, err := ioutil.ReadFile(CLIENTB_DEFAULT_FOLDER + fileName)

	server.Process.Kill()
	clientB.Process.Kill()
	clientA.Process.Kill()

	if !reflect.DeepEqual(actualFileContent, expectedFileContent) || err != nil {
		t.Fatalf("Expected: %d, got: %d \nError: %v", len(expectedFileContent), len(actualFileContent), err)
	}

	cleanup(CLIENTB_DEFAULT_FOLDER + fileName)
}

func setupTest() (filePath string, fileName string, fileContent []byte) {
	fileName = "autoTest.txt"
	filePath = "C:\\Users\\dejim\\source\\repos\\file-transfer-go\\tests\\scenarios\\one-to-one-transfer\\senderFiles\\"
	fileContent = []byte("This is a txt file for testing the one to one transfer!!!\n")
	fileToRemove := fmt.Sprint(CLIENTB_DEFAULT_FOLDER, fileName)
	os.Remove(fileToRemove)
	os.Remove(filePath + fileName)

	file, _ := os.Create(filePath + fileName)
	defer file.Close()

	file.Write(fileContent)
	return
}

func cleanup(filePath string) {
	os.Remove(filePath)
}
