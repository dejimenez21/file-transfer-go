package main

type request struct {
	Method      string
	Meta        metaData
	Channels    []string
	FileInfo    file
	FileContent []byte
}

func (req request) getMessageType() string {
	return "request"
}

type metaData struct {
	HasFileContent bool
	SenderAddress  string
}

type delivery struct {
	DeliveryId int64
	Seq        int64
	Size       int
	Content    []byte
}

func (del delivery) getMessageType() string {
	return "chunk"
}

type cftpMessage interface {
	getMessageType() string
}
