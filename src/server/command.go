package main

type command struct {
	Method   string
	Meta     metaData
	Channels []string
	FileInfo file
}

type metaData struct {
	SenderAddress string
	RequestId     int
}
