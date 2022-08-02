package main

type receiveCmd struct {
	channels   []string
	folderPath string
}

type sendCmd struct {
	channels []string
	filePath string
}
