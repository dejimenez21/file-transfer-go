package main

type request struct {
	Method   string
	Meta     metaData
	Channels []string
	FileInfo file
}

type metaData struct {
	HasFileContent bool
	SenderAddress  string
}
